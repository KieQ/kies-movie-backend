package db

import (
	"context"
	"errors"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"kies-movie-backend/model/table"
	"kies-movie-backend/utils"
	"strings"
	"time"
)

func AddVideo(ctx context.Context, video *table.Video) error {
	logs.CtxInfo(ctx, "added video=%v", utils.ToJSON(video))
	if video == nil {
		return errors.New("user is nil")
	}
	err := db.Table(table.NameVideo).Create(video).Error
	return err
}

func UpdateVideoByID(ctx context.Context, account string, id int64, updateData map[string]interface{}) (int64, error) {
	logs.CtxInfo(ctx, "account=%v, update data=%v, id=%v", account, utils.ToJSON(updateData), id)
	if updateData == nil {
		return 0, errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	result := db.Table(table.NameVideo).Where("id = ? AND user_account = ?", id, account).Updates(updateData)
	return result.RowsAffected, result.Error
}

func UpdateVideosByIDs(ctx context.Context, account string, ids []int64, updateData map[string]interface{}) (int64, error) {
	logs.CtxInfo(ctx, "account=%v, update data=%v, ids=%v", account, utils.ToJSON(updateData), utils.ToJSON(ids))
	if updateData == nil {
		return 0, errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	result := db.Table(table.NameVideo).Where("id in ? AND user_account = ?", ids, account).Updates(updateData)
	return result.RowsAffected, result.Error
}

func GetVideo(ctx context.Context, where map[string]interface{}) ([]*table.Video, error) {
	logs.CtxInfo(ctx, "where condition=%v", utils.ToJSON(where))
	result := make([]*table.Video, 0, 1)
	err := db.Table(table.NameVideo).Where(where).Find(&result).Error
	return result, err
}

func GetVideoByID(ctx context.Context, id int64) (*table.Video, error) {
	logs.CtxInfo(ctx, "id=%v", utils.ToJSON(id))
	video := new(table.Video)
	err := db.Table(table.NameVideo).Where("id = ?", id).Last(&video).Error
	return video, err
}

func GetVideoCount(ctx context.Context, where map[string]interface{}) (int64, error) {
	logs.CtxInfo(ctx, "where condition=%v", utils.ToJSON(where))
	var count int64
	err := db.Table(table.NameVideo).Where(where).Count(&count).Error
	return count, err
}

func GetVideoCountByType(ctx context.Context, types []table.VideoType) (int64, error) {
	logs.CtxInfo(ctx, "types=%v", utils.ToJSON(types))
	var count int64
	err := db.Table(table.NameVideo).Where("video_type in ?", types).Count(&count).Error
	return count, err
}

func GetOrderedVideoWithLimit(ctx context.Context, where map[string]interface{}, offset, limit int, orderBy string, desc bool) ([]*table.Video, error) {
	logs.CtxInfo(ctx, "where condition=%v, offset=%v, limit=%v", where, offset, limit)
	result := make([]*table.Video, 0, 1)
	if desc && !strings.HasSuffix(strings.TrimSpace(strings.ToLower(orderBy)), "desc") {
		orderBy += " DESC"
	}
	err := db.Table(table.NameVideo).Where(where).Offset(offset).Limit(limit).Order(orderBy).Find(&result).Error
	return result, err
}

func GetVideoByType(ctx context.Context, types []int, where map[string]interface{}) ([]*table.Video, error) {
	logs.CtxInfo(ctx, "where condition=%v, types=%v", utils.ToJSON(where), utils.ToJSON(types))
	result := make([]*table.Video, 0, 1)
	err := db.Table(table.NameVideo).Where("video_type in ?", types).Where(where).Find(&result).Error
	return result, err
}

func GetVideoByTypeWithLimitation(ctx context.Context, types []table.VideoType, where map[string]interface{}, offset, limit int, orderBy string, desc bool) ([]*table.Video, error) {
	logs.CtxInfo(ctx, "where condition=%v, types=%v, offset=%v, limit=%v, orderBy=%v, desc=%v", utils.ToJSON(where), utils.ToJSON(types), offset, limit, orderBy, desc)
	result := make([]*table.Video, 0, 1)
	if desc && !strings.HasSuffix(strings.TrimSpace(strings.ToLower(orderBy)), "desc") {
		orderBy += " DESC"
	}
	err := db.Table(table.NameVideo).Where("video_type in ?", types).Where(where).Offset(offset).Limit(limit).Order(orderBy).Find(&result).Error
	return result, err
}

func DeleteVideoByID(ctx context.Context, account string, id int64) (int64, error) {
	logs.CtxInfo(ctx, "account=%v, id=%v", account, id)
	result := db.Table(table.NameVideo).Where("user_account=? AND id=?", account, id).Delete(nil)
	return result.RowsAffected, result.Error
}
