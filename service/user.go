package service

import (
	"kies-movie-backend/dto"
	"kies-movie-backend/model/table"
	"kies-movie-backend/utils"
)

func TransUserTableToDTO(user *table.User) *dto.User {
	if user == nil {
		return nil
	}
	return &dto.User{
		Account:          user.Account,
		NickName:         user.NickName,
		Profile:          user.Profile,
		Phone:            user.Phone,
		Email:            user.Email,
		Gender:           user.Gender,
		SelfIntroduction: user.SelfIntroduction,
		PreferTags:       utils.FromJSON[[]string](user.PreferTags),
		CreateTime:       user.CreateTime.Unix(),
	}
}
