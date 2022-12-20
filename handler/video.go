package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"kies-movie-backend/constant"
	"kies-movie-backend/dto"
	"kies-movie-backend/i18n"
	"kies-movie-backend/model/db"
	"kies-movie-backend/model/table"
	"kies-movie-backend/service"
	"kies-movie-backend/utils"
	"strconv"
	"time"
)

func VideoList(c *gin.Context) {
	page := c.Query("page")
	size := c.Query("size")
	logs.CtxInfo(c, "page=%v, size=%v", page, size)

	pageVal, err := strconv.ParseInt(page, 10, 32)
	if err != nil {
		logs.CtxWarn(c, "failed to parse page, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}

	sizeVal, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		logs.CtxWarn(c, "failed to parse size, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}

	account := c.GetString(constant.Account)
	var videos []*table.Video
	var count int64
	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		videos, err = db.GetOrderedVideoWithLimit(c, map[string]interface{}{"user_account": account}, int(pageVal)*int(sizeVal), int(sizeVal), "id", true)
		if err != nil {
			logs.CtxWarn(c, "failed to read videos, err=%v", err)
			return err
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		count, err = db.GetVideoCount(c, map[string]interface{}{"user_account": account})
		if err != nil {
			logs.CtxWarn(c, "failed to read videos count, err=%v", err)
			return err
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}

	items := service.TransForVideoListDTO(videos)
	resp := &dto.VideoListResponse{
		Total: count,
		Page:  int32(pageVal),
		Size:  int32(pageVal),
		Items: items,
	}

	logs.CtxInfo(c, "【KIE DEBUG】%v", utils.ToJSON(items))
	OnSuccess(c, resp)
}

func VideoDetail(c *gin.Context) {
	OnSuccess(c, nil)
}

func VideoLike(c *gin.Context) {
	req := dto.VideoLikeRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))

	account := c.GetString(constant.Account)

	rows, err := db.UpdateVideoByID(c, account, req.ID, map[string]interface{}{
		"liked":       req.Liked,
		"update_time": time.Now(),
	})
	if err != nil {
		logs.CtxWarn(c, "failed to update %v", req.ID)
		OnFail(c, constant.FailedToProcess)
		return
	}
	if rows == 0 {
		logs.CtxWarn(c, "video with id %v does not belong to account %v", req.ID, account)
		OnFail(c, constant.NoAuthority)
		return
	}
	OnSuccess(c, nil)
}

func VideoUpdate(c *gin.Context) {
	OnSuccess(c, nil)
}

func VideoAdd(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	posterPath := c.PostForm("poster_path")
	region := c.PostForm("region")
	videoType := c.PostForm("video_type")
	link := c.PostForm("link")
	linkType := c.PostForm("link_type")
	backdropPath := c.PostForm("backdrop_path")

	if name == "" || posterPath == "" || region == "" || linkType == "" || videoType == "" {
		logs.CtxWarn(c, "needed data is missing, name=%v, posterPath=%v, region=%v, linkType=%v, videoType=%v", name, posterPath, region, linkType, videoType)
		OnFail(c, constant.RequestParameterError)
		return
	}

	account := c.GetString(constant.Account)

	//For torrent file, translate it to magnet and store in the db
	if linkType == "3" {
		file, err := c.FormFile("link_file")
		if err != nil {
			logs.CtxWarn(c, "failed to get file, err=%v", err)
			OnFail(c, constant.RequestParameterError)
			return
		}
		reader, err := file.Open()
		if err != nil {
			logs.CtxWarn(c, "failed to open file, err=%v", err)
			OnFail(c, constant.RequestParameterError)
			return
		}
		mi, err := metainfo.Load(reader)
		if err != nil {
			logs.CtxWarn(c, "failed to load meta info from file, err=%v", err)
			OnFail(c, constant.FailedToProcess)
			return
		}
		link = mi.Magnet(nil, nil).String()
		linkType = strconv.FormatInt(int64(table.LinkTypeMagnet), 10)
	}

	//Video Type Check
	var videoTypeValue table.VideoType
	switch videoType {
	case "0":
		videoTypeValue = table.VideoTypeMovie
	case "1":
		videoTypeValue = table.VideoTypeTV
	case "2":
		videoTypeValue = table.VideoTypeMoviePrivate
	case "3":
		videoTypeValue = table.VideoTypeTVPrivate
	default:
		logs.CtxWarn(c, "unknown video type, videoType=%v", videoType)
		OnFail(c, constant.RequestParameterError)
		return
	}

	//Region Check
	if !utils.Contain([]string{"en", "zh", "jp"}, region) {
		logs.CtxWarn(c, "unknown region, region=%v", region)
		OnFail(c, constant.RequestParameterError)
		return
	}

	//LinkType Check
	linkTypeValueInt, err := strconv.ParseInt(linkType, 10, 64)
	if err != nil {
		logs.CtxWarn(c, "failed to parse linkType as int, linkType=%v", linkType)
		OnFail(c, constant.RequestParameterError)
		return
	}
	linkTypeValue := table.LinkType(linkTypeValueInt)
	if !utils.Contain([]table.LinkType{table.LinkTypeLinkAddress, table.LinkTypeNoLink, table.LinkTypeMagnet}, linkTypeValue) {
		logs.CtxWarn(c, "unknown linkType Value, linkTypeValue=%v", linkTypeValue)
		OnFail(c, constant.RequestParameterError)
		return
	}

	err = db.AddVideo(c, &table.Video{
		VideoName:        name,
		VideoDescription: description,
		VideoSize:        0,
		VideoType:        videoTypeValue,
		Region:           region,
		Link:             link,
		LinkType:         linkTypeValue,
		Location:         "",
		PosterPath:       posterPath,
		BackdropPath:     backdropPath,
		UserAccount:      account,
		Tags:             "",
		Liked:            false,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	})
	if err != nil {
		logs.CtxWarn(c, "failed to add video to database, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToAddVideo)
		return
	}

	OnSuccess(c, nil)
}

func VideoDelete(c *gin.Context) {
	req := dto.VideoDeleteRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))

	account := c.GetString(constant.Account)
	rows, err := db.DeleteVideoByID(c, account, req.ID)
	if err != nil {
		logs.CtxWarn(c, "failed to delete %v", req.ID)
		OnFail(c, constant.FailedToProcess)
		return
	}
	if rows == 0 {
		logs.CtxWarn(c, "video with id %v does not belong to account %v", req.ID, account)
		OnFail(c, constant.NoAuthority)
		return
	}
	OnSuccess(c, nil)
}

func VideoDownload(c *gin.Context) {
	OnSuccess(c, nil)
}

func VideoClone(c *gin.Context) {
	req := dto.VideoCloneRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))

	account := c.GetString(constant.Account)

	videos, err := db.GetVideo(c, map[string]interface{}{"id": req.ID})
	if err != nil {
		logs.CtxWarn(c, "failed to get video, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}

	if len(videos) == 0 {
		logs.CtxWarn(c, "can't find video, it might be deleted")
		OnFailWithMessage(c, constant.FailedToProcess, i18n.VideoMightBeDeleted)
		return
	}

	video := videos[0]
	if video.UserAccount == account {
		logs.CtxWarn(c, "account and user account are the same")
		OnFailWithMessage(c, constant.FailedToProcess, i18n.CannotCloneYourOwnMovie)
		return
	}

	err = db.AddVideo(c, &table.Video{
		VideoName:        video.VideoName,
		VideoDescription: video.VideoDescription,
		VideoSize:        video.VideoSize,
		VideoType:        video.VideoType,
		Region:           video.Region,
		Link:             video.Link,
		LinkType:         video.LinkType,
		Location:         video.Location,
		PosterPath:       video.PosterPath,
		BackdropPath:     video.BackdropPath,
		UserAccount:      account,
		Tags:             video.Tags,
		Liked:            false,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	})
	if err != nil {
		logs.CtxWarn(c, "failed to add video, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToAddVideo)
		return
	}
	OnSuccess(c, nil)
}
