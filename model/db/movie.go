package db

import (
	"context"
	"errors"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"kies-movie-backend/model/table"
	"kies-movie-backend/utils"
	"time"
)

func AddMovie(ctx context.Context, movie *table.Movie) error {
	logs.CtxInfo(ctx, "added movie=%v", utils.ToJSON(movie))
	if movie == nil {
		return errors.New("user is nil")
	}
	err := movieDB.Table(movie.Table()).Create(movie).Error
	return err
}

func UpdateMovie(ctx context.Context, name string, updateData map[string]interface{}) error {
	logs.CtxInfo(ctx, "update data=%v, name=%v", utils.ToJSON(updateData), name)
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := movieDB.Table(new(table.User).Table()).Where("name = ?", name).Updates(updateData).Error
	return err
}

func UpdateMovies(ctx context.Context, names []string, updateData map[string]interface{}) error {
	logs.CtxInfo(ctx, "update data=%v, names=%v", utils.ToJSON(updateData), utils.ToJSON(names))
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := movieDB.Table(new(table.User).Table()).Where("name in ?", names).Updates(updateData).Error
	return err
}

func GetMovie(ctx context.Context, where map[string]interface{}) ([]*table.Movie, error) {
	logs.CtxInfo(ctx, "where condition=%v", utils.ToJSON(where))
	result := make([]*table.Movie, 0, 1)
	err := movieDB.Table(new(table.Movie).Table()).Where(where).Find(&result).Error
	return result, err
}
