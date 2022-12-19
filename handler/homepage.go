package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-movie-backend/constant"
	"kies-movie-backend/dto"
	"kies-movie-backend/i18n"
	"kies-movie-backend/model/db"
	"kies-movie-backend/model/table"
	"math/rand"
	"strconv"
)

func HomepageContent(c *gin.Context) {
	lang := c.GetString(constant.Language)
	where := make(map[string]interface{})
	if len(lang) > 0 {
		where["region"] = lang
	}

	allVideo, err := db.GetVideoByType(c, []int{int(table.VideoTypeMovie), int(table.VideoTypeTV)}, where)
	if err != nil {
		logs.CtxWarn(c, "failed to get video, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToFindMovieOrTV)
		return
	}

	result := &dto.HomepageContentResponse{CarouselItems: make([]*dto.CarouselItem, 0, 3)}
	for _, v := range rand.Perm(len(allVideo)) {
		result.CarouselItems = append(result.CarouselItems, &dto.CarouselItem{
			PosterImage:     allVideo[v].PosterPath,
			BackgroundImage: allVideo[v].BackdropPath,
			Title:           allVideo[v].VideoName,
			Content:         allVideo[v].VideoDescription,
			Value:           strconv.FormatInt(allVideo[v].ID, 10),
		})
		if len(result.CarouselItems) == 3 {
			break
		}
	}

	OnSuccess(c, result)
}
