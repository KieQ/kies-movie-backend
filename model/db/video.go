package db

import (
	"context"
	"errors"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"kies-movie-backend/model/table"
	"kies-movie-backend/utils"
	"time"
)

func AddVideo(ctx context.Context, video *table.Video) error {
	logs.CtxInfo(ctx, "added video=%v", utils.ToJSON(video))
	if video == nil {
		return errors.New("user is nil")
	}
	err := db.Table(video.Table()).Create(video).Error
	return err
}

func UpdateVideo(ctx context.Context, name string, updateData map[string]interface{}) error {
	logs.CtxInfo(ctx, "update data=%v, name=%v", utils.ToJSON(updateData), name)
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := db.Table(new(table.Video).Table()).Where("name = ?", name).Updates(updateData).Error
	return err
}

func UpdateVideos(ctx context.Context, names []string, updateData map[string]interface{}) error {
	logs.CtxInfo(ctx, "update data=%v, names=%v", utils.ToJSON(updateData), utils.ToJSON(names))
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := db.Table(new(table.Video).Table()).Where("name in ?", names).Updates(updateData).Error
	return err
}

func GetVideo(ctx context.Context, where map[string]interface{}) ([]*table.Video, error) {
	logs.CtxInfo(ctx, "where condition=%v", utils.ToJSON(where))
	result := make([]*table.Video, 0, 1)
	err := db.Table(new(table.Video).Table()).Where(where).Find(&result).Error
	return result, err
}

func GetVideoByType(ctx context.Context, types []int, where map[string]interface{}) ([]*table.Video, error) {
	logs.CtxInfo(ctx, "where condition=%v, types=%v", utils.ToJSON(where), utils.ToJSON(types))
	result := make([]*table.Video, 0, 1)
	err := db.Table(new(table.Video).Table()).Where("video_type in %v", types).Where(where).Find(&result).Error
	return result, err
}
