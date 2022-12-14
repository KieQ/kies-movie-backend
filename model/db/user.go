package db

import (
	"context"
	"errors"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"kies-movie-backend/model/table"
	"kies-movie-backend/utils"
	"time"
)

func AddUser(ctx context.Context, user *table.User) error {
	logs.CtxInfo(ctx, "added user=%v", utils.ToJSON(user))
	if user == nil {
		return errors.New("user is nil")
	}
	err := db.Table(table.NameUser).Create(user).Error
	return err
}

func UpdateUser(ctx context.Context, account string, updateData map[string]interface{}) error {
	logs.CtxInfo(ctx, "update data=%v, account=%v", utils.ToJSON(updateData), account)
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := db.Table(table.NameUser).Where("account = ?", account).Updates(updateData).Error
	return err
}

func UpdateUsers(ctx context.Context, accounts []string, updateData map[string]interface{}) error {
	logs.CtxInfo(ctx, "update data=%v, accounts=%v", utils.ToJSON(updateData), utils.ToJSON(accounts))
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := db.Table(table.NameUser).Where("account in ?", accounts).Updates(updateData).Error
	return err
}

func GetUser(ctx context.Context, where map[string]interface{}) ([]*table.User, error) {
	logs.CtxInfo(ctx, "where condition=%v", utils.ToJSON(where))
	result := make([]*table.User, 0, 1)
	err := db.Table(table.NameUser).Where(where).Find(&result).Error
	return result, err
}

func GetUsersWithAccounts(ctx context.Context, accounts []string) ([]*table.User, error) {
	logs.CtxInfo(ctx, "accounts=%v", utils.ToJSON(accounts))
	result := make([]*table.User, 0, len(accounts))
	err := db.Table(table.NameUser).Where("account in ?", accounts).Find(&result).Error
	return result, err
}
