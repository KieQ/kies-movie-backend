package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-movie-backend/constant"
	"kies-movie-backend/dto"
	"kies-movie-backend/model/db"
	"kies-movie-backend/model/table"
	"kies-movie-backend/service"
	"kies-movie-backend/utils"
	"time"
)

func UserAdd(c *gin.Context) {
	req := dto.UserAddRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))
	user := &table.User{
		Account:          req.Account,
		Password:         req.Password,
		NickName:         req.NickName,
		Profile:          req.Profile,
		Phone:            req.Phone,
		Email:            req.Email,
		Gender:           req.Gender,
		SelfIntroduction: req.SelfIntroduction,
		PreferTags:       utils.ToJSON(req.PreferTags),
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
	err = db.AddUser(c, user)
	if err != nil {
		logs.CtxWarn(c, "failed to add user %v, err=%v", req.Account, err)
		OnFail(c, constant.ServiceError)
		return
	}
	OnSuccess(c, nil)
}

func UserUpdate(c *gin.Context) {
	req := dto.UserUpdateRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	if req.Account == "" {
		logs.CtxWarn(c, "account is empty")
		OnFail(c, constant.RequestParameterError)
		return
	}

	updateData := make(map[string]interface{})
	utils.AddToMap(updateData, req.Password, "password")
	utils.AddToMap(updateData, req.Email, "email")
	utils.AddToMap(updateData, req.Phone, "phone")
	utils.AddToMap(updateData, req.Gender, "gender")
	utils.AddToMap(updateData, req.SelfIntroduction, "self_introduction")
	utils.AddToMap(updateData, req.Profile, "profile")
	utils.AddToMap(updateData, req.NickName, "nick_name")
	if len(req.PreferTags) != 0 {
		updateData["prefer_tags"] = utils.ToJSON(req.PreferTags)
	}

	err = db.UpdateUser(c, req.Account, updateData)
	if err != nil {
		logs.CtxWarn(c, "failed to update user %v, err=%v", req.Account, err)
		OnFail(c, constant.ServiceError)
		return
	}
	OnSuccess(c, nil)
}

func UserDetail(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		logs.CtxWarn(c, "account is empty")
		OnFail(c, constant.RequestParameterError)
		return
	}

	user, err := db.GetUser(c, map[string]interface{}{"account": account})
	if err != nil {
		logs.CtxWarn(c, "failed to get user, err=%v", err)
		OnFail(c, constant.ServiceError)
		return
	}
	if len(user) == 0 {
		logs.CtxWarn(c, "length of user is empty")
		OnFailWithMessage(c, constant.RequestParameterError, "there is no user found")
		return
	}

	OnSuccess(c, dto.UserDetailResponse{
		User: service.TransUserTableToDTO(user[0]),
	})
}

func UserList(c *gin.Context) {
	user, err := db.GetUser(c, nil)
	if err != nil {
		logs.CtxWarn(c, "failed to get user, err=%v", err)
		OnFail(c, constant.ServiceError)
		return
	}

	result := dto.UserListResponse{
		Users: make([]*dto.User, 0, 10),
	}
	for _, item := range user {
		result.Users = append(result.Users, service.TransUserTableToDTO(item))
	}

	OnSuccess(c, result)
}
