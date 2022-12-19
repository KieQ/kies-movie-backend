package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"kies-movie-backend/constant"
	"kies-movie-backend/dto"
	"kies-movie-backend/i18n"
	"kies-movie-backend/model/db"
	"kies-movie-backend/model/table"
	"kies-movie-backend/service"
	"strconv"
)

func NotLoginVideoList(c *gin.Context) {
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

	var videos []*table.Video
	var count int64
	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		videos, err = db.GetVideoByTypeWithLimitation(c, []table.VideoType{table.VideoTypeTV, table.VideoTypeMovie}, nil, int(pageVal)*int(sizeVal), int(sizeVal), "id", true)
		if err != nil {
			logs.CtxWarn(c, "failed to read videos, err=%v", err)
			return err
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		count, err = db.GetVideoCountByType(c, []table.VideoType{table.VideoTypeTV, table.VideoTypeMovie})
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

	accounts := make([]string, 0, len(videos))
	for _, item := range videos {
		accounts = append(accounts, item.UserAccount)
	}

	users, err := db.GetUsersWithAccounts(c, accounts)
	if err != nil {
		logs.CtxWarn(c, "failed to get users, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindUsers)
		return
	}

	items := service.TransForNotLoginVideoListDTO(videos, users)
	resp := &dto.NotLoginVideoListResponse{
		Total: count,
		Page:  int32(pageVal),
		Size:  int32(sizeVal),
		Items: items,
	}
	OnSuccess(c, resp)
}

func NotLoginVideoDetail(c *gin.Context) {
	OnSuccess(c, nil)
}
