package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"kies-movie-backend/api/tmdb"
	"kies-movie-backend/constant"
	"kies-movie-backend/dto"
	"kies-movie-backend/i18n"
	"kies-movie-backend/utils"
	"math/rand"
	"strings"
	"time"
)

func HomepageContent(c *gin.Context) {
	var lang, originalLang string
	switch strings.ToLower(c.GetString(i18n.ContextLanguage)) {
	case "", "en":
		lang, originalLang = "en-US", "en"
	case "zh-cn":
		lang, originalLang = "zh-CN", "zh"
	}
	year := time.Now().AddDate(0, -6, 0).Year() - rand.Intn(15)
	var movie []map[string]interface{}
	var tv []map[string]interface{}
	eg := errgroup.Group{}
	eg.Go(func() error {
		temp, err := tmdb.DiscoverMovie(c, lang, originalLang, year)
		for _, item := range temp {
			if utils.DowncastWithDefault[string](item["poster_path"], "") != "" &&
				utils.DowncastWithDefault[string](item["backdrop_path"], "") != "" {
				movie = append(movie, item)
			}
		}
		return err
	})
	eg.Go(func() error {
		temp, err := tmdb.DiscoverTV(c, lang, originalLang, year)
		for _, item := range temp {
			if utils.DowncastWithDefault[string](item["poster_path"], "") != "" &&
				utils.DowncastWithDefault[string](item["backdrop_path"], "") != "" {
				tv = append(tv, item)
			}
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		logs.CtxWarn(c, "failed to discover movie or tv, err=%v", err)
		OnFail(c, constant.ServiceError)
		return
	}

	if len(movie)+len(tv) == 0 {
		logs.CtxWarn(c, "no movie and tv fetched")
		OnFailWithMessage(c, constant.ServiceError, i18n.NoMovieOrTVFound)
		return
	}

	result := &dto.HomepageContentResponse{CarouselItems: make([]*dto.CarouselItem, 0, 3)}
	for _, item := range utils.Sample(movie, 2) {
		result.CarouselItems = append(result.CarouselItems, &dto.CarouselItem{
			PosterImage:     tmdb.WrapImage(utils.DowncastWithDefault[string](item["poster_path"], "")),
			BackgroundImage: tmdb.WrapImage(utils.DowncastWithDefault[string](item["backdrop_path"], "")),
			Title:           utils.DowncastWithDefault[string](item["title"], ""),
			Content:         utils.DowncastWithDefault[string](item["overview"], ""),
			Value:           "",
		})
	}

	for _, item := range utils.Sample(tv, 1) {
		result.CarouselItems = append(result.CarouselItems, &dto.CarouselItem{
			PosterImage:     tmdb.WrapImage(utils.DowncastWithDefault[string](item["poster_path"], "")),
			BackgroundImage: tmdb.WrapImage(utils.DowncastWithDefault[string](item["backdrop_path"], "")),
			Title:           utils.DowncastWithDefault[string](item["name"], ""),
			Content:         utils.DowncastWithDefault[string](item["overview"], ""),
			Value:           "",
		})
	}

	OnSuccess(c, result)
}
