package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"kies-movie-backend/utils"
	"os"
)

var appKey string

func init() {
	appKey = os.Getenv("TMDB_APP_KEY")
	if appKey == "" {
		panic("TMDB_APP_KEY is needed")
	}
}

type response struct {
	Results []map[string]interface{} `json:"results"`
}

func SearchMulti(ctx context.Context, name string, lang string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/multi?api_key=%v&language=%v&page=1&include_adult=false&query=%v", appKey, lang, name)
	logs.CtxInfo(ctx, "start get %v", url)
	content, err := utils.HTTPGet(url)
	if err != nil {
		return nil, err
	}
	result := new(response)
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func DiscoverTV(ctx context.Context, lang, originalLang string, year int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/discover/tv?api_key=%v&language=%v&with_original_language=%v&first_air_date_year=%v", appKey, lang, originalLang, year)
	logs.CtxInfo(ctx, "start get %v", url)
	content, err := utils.HTTPGet(url)
	if err != nil {
		return nil, err
	}
	result := new(response)
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func DiscoverMovie(ctx context.Context, lang, originalLang string, year int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?api_key=%v&language=%v&with_original_language=%v&year=%v", appKey, lang, originalLang, year)
	logs.CtxInfo(ctx, "start get %v", url)
	content, err := utils.HTTPGet(url)
	if err != nil {
		return nil, err
	}
	result := new(response)
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func WrapImage(image string) string {
	return "https://image.tmdb.org/t/p/original" + image
}
