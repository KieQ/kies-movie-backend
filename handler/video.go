package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"kies-movie-backend/constant"
	"kies-movie-backend/download"
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
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logs.CtxWarn(c, "failed to parse idStr %v", idStr)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "id=%v", id)

	account := c.GetString(constant.Account)
	video, err := db.GetVideoByID(c, id)
	if err != nil {
		logs.CtxWarn(c, "failed to fetch video, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}
	if video.UserAccount != account {
		logs.CtxWarn(c, "user account is not the same for id %v, account in context=%v, account of video=%v", id, account, video.UserAccount)
		OnFail(c, constant.NoAuthority)
		return
	}

	resp := &dto.VideoDetailResponse{
		ID:           video.ID,
		Name:         video.VideoName,
		Description:  video.VideoDescription,
		VideoType:    video.VideoType,
		Region:       video.Region,
		Link:         video.Link,
		Files:        video.Files,
		Downloaded:   video.Downloaded,
		PosterPath:   video.PosterPath,
		BackdropPath: video.BackdropPath,
		UserAccount:  video.UserAccount,
		Tags:         video.Tags,
		Liked:        video.Liked,
	}

	switch video.LinkType {
	case table.LinkTypeNoLink:
		resp.LinkType = "0"
	case table.LinkTypeLinkAddress:
		resp.LinkType = "1"
	case table.LinkTypeMagnet:
		resp.LinkType = "2"
	}

	OnSuccess(c, resp)
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
	idStr := c.PostForm("id")
	name := c.PostForm("name")
	description := c.PostForm("description")
	posterPath := c.PostForm("poster_path")
	region := c.PostForm("region")
	videoType := c.PostForm("video_type")
	link := c.PostForm("link")
	linkType := c.PostForm("link_type")
	backdropPath := c.PostForm("backdrop_path")

	if idStr == "" || name == "" || posterPath == "" || region == "" || linkType == "" || videoType == "" {
		logs.CtxWarn(c, "needed data is missing, id=%v, name=%v, posterPath=%v, region=%v, linkType=%v, videoType=%v", idStr, name, posterPath, region, linkType, videoType)
		OnFail(c, constant.RequestParameterError)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logs.CtxWarn(c, "failed to parse idStr %v, err=%v", idStr, err)
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
		linkType = "2"
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
	var linkTypeValue = table.LinkTypeNoLink
	switch linkType {
	case "0":
		linkTypeValue = table.LinkTypeNoLink
	case "1":
		linkTypeValue = table.LinkTypeLinkAddress
	case "2", "3":
		linkTypeValue = table.LinkTypeMagnet
	default:
		logs.CtxWarn(c, "unknown linkType Value, linkType=%v", linkType)
		OnFail(c, constant.RequestParameterError)
		return
	}

	rows, err := db.UpdateVideoByID(c, account, id, map[string]interface{}{
		"video_name":        name,
		"video_description": description,
		"video_type":        videoTypeValue,
		"region":            region,
		"link":              link,
		"link_type":         linkTypeValue,
		"poster_path":       posterPath,
		"backdrop_path":     backdropPath,
		"update_time":       time.Now(),
	})
	if err != nil {
		logs.CtxWarn(c, "failed to update video to database, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToAddVideo)
		return
	}
	if rows == 0 {
		logs.CtxInfo(c, "movie might be removed")
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}

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
		linkType = "2"
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
	var linkTypeValue = table.LinkTypeNoLink
	switch linkType {
	case "0":
		linkTypeValue = table.LinkTypeNoLink
	case "1":
		linkTypeValue = table.LinkTypeLinkAddress
	case "2", "3":
		linkTypeValue = table.LinkTypeMagnet
	default:
		logs.CtxWarn(c, "unknown linkType Value, linkType=%v", linkType)
		OnFail(c, constant.RequestParameterError)
		return
	}

	err := db.AddVideo(c, &table.Video{
		VideoName:        name,
		VideoDescription: description,
		VideoType:        videoTypeValue,
		Region:           region,
		Link:             link,
		LinkType:         linkTypeValue,
		Files:            "",
		Downloaded:       false,
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

	video, err := db.GetVideoByID(c, req.ID)
	if err != nil {
		logs.CtxWarn(c, "failed to fetch video, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}
	if video.UserAccount != account {
		logs.CtxWarn(c, "user account is not the same for id %v, account in context=%v, account of video=%v", req.ID, account, video.UserAccount)
		OnFail(c, constant.NoAuthority)
		return
	}

	//Delete from DB
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

	//Remove from downloadingMap
	mag, err := metainfo.ParseMagnetUri(video.Link)
	if err == nil {
		if v, exist, _ := download.GetFromDownloadingMap(mag.InfoHash.HexString()); exist {
			v.Torrent.Drop()
			download.RemoveFromDownloadingMap(mag.InfoHash.HexString())
		}
	}

	//Delete downloaded files
	if video.Files != "" {
		go func() {
			filenames := utils.FromJSON[[]string](video.Files)
			for i := 0; i < 10; i++ {
				filenames = download.DeleteWholeDirectory(c, filenames)
				if len(filenames) == 0 {
					break
				}
				time.Sleep(time.Second)
			}
		}()
	}

	OnSuccess(c, nil)
}

func VideoAvailableFiles(c *gin.Context) {
	idStr, ok := c.GetQuery("id")
	if !ok {
		logs.CtxWarn(c, "id does not exist")
		OnFail(c, constant.RequestParameterError)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logs.CtxWarn(c, "id is not integer, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}

	account := c.GetString(constant.Account)

	video, err := db.GetVideoByID(c, id)
	if err != nil {
		logs.CtxWarn(c, "failed to fetch video with id %v, err=%v", id, err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}

	if video.UserAccount != account {
		logs.CtxWarn(c, "user account is not the same for id %v, account in context=%v, account of video=%v", id, account, video.UserAccount)
		OnFail(c, constant.NoAuthority)
		return
	}

	resp := dto.VideoAvailableFilesResponse{}
	if video.LinkType == table.LinkTypeMagnet {
		infoHash, files, timeout, err := download.ShowFilesInMagnet(c, video.Files, video.Link)
		if err != nil {
			logs.CtxWarn(c, "failed to retrieve files information with id %v", id)
			OnFail(c, constant.FailedToProcess)
			return
		}

		resp.InfoHash = infoHash
		resp.Timeout = timeout
		respFiles := make([]dto.VideoAvailableFileInfo, 0, len(files))
		for _, item := range files {
			respFiles = append(respFiles, dto.VideoAvailableFileInfo{
				Path:            item.Path(),
				DisplayPath:     item.DisplayPath(),
				DownloadedBytes: item.BytesCompleted(),
				TotalBytes:      item.Length(),
				Downloading:     item.Priority() != torrent.PiecePriorityNone,
			})
		}
		resp.Files = respFiles
	} else {
		logs.CtxWarn(c, "no link should not call this method, id=%v", id)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.NoLinkCannotBeProcessed)
		return
	}
	OnSuccess(c, resp)
}

func VideoDownload(c *gin.Context) {
	req := dto.VideoDownloadRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))

	if len(req.Files) == 0 || len(req.InfoHash) == 0 {
		logs.CtxWarn(c, "Files and InfoHash cannot be empty")
		OnFail(c, constant.RequestParameterError)
		return
	}

	account := c.GetString(constant.Account)

	video, err := db.GetVideoByID(c, req.ID)
	if err != nil {
		logs.CtxWarn(c, "failed to fetch video, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}

	if video.UserAccount != account {
		logs.CtxWarn(c, "no authority, account=%v, user_account=%v", account, video.UserAccount)
		OnFail(c, constant.NoAuthority)
		return
	}

	files, exist, err := download.StartDownloadSelectFileAsync(c, req.ID, account, req.InfoHash, req.Files)
	if err != nil {
		logs.CtxWarn(c, "failed to start download, err=%v", err)
		OnFail(c, constant.FailedToProcess)
		return
	}
	if !exist {
		logs.CtxWarn(c, "infoHash %v is not added to downloadingMap", req.InfoHash)
		OnFail(c, constant.FailedToProcess)
		return
	}

	filenames := make([]string, 0, len(files))
	for _, item := range files {
		filenames = append(filenames, item.Path())
	}
	rows, err := db.UpdateVideoByID(c, account, req.ID, map[string]interface{}{"files": utils.ToJSON(filenames), "downloaded": false})
	if err != nil {
		logs.CtxWarn(c, "failed to update files, err=%v", err)
		OnFail(c, constant.FailedToProcess)
		return
	}
	if rows == 0 {
		logs.CtxInfo(c, "movie might be removed")
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}

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

	video, err := db.GetVideoByID(c, req.ID)
	if err != nil {
		logs.CtxWarn(c, "failed to get video, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}
	if video.UserAccount == account {
		logs.CtxWarn(c, "account and user account are the same")
		OnFailWithMessage(c, constant.FailedToProcess, i18n.CannotCloneYourOwnMovie)
		return
	}

	err = db.AddVideo(c, &table.Video{
		VideoName:        video.VideoName,
		VideoDescription: video.VideoDescription,
		VideoType:        video.VideoType,
		Region:           video.Region,
		Link:             video.Link,
		LinkType:         video.LinkType,
		Files:            video.Files,
		Downloaded:       video.Downloaded,
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
