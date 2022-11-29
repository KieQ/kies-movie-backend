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

func SessionManageLogin(c *gin.Context) {
	req := dto.SessionManageLoginRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))

	// Parameter check
	if req.Account == "" || req.Password == "" {
		logs.CtxWarn(c, "required parameters are missing")
		OnFail(c, constant.RequestParameterError)
		return
	}

	// Existence check
	users, err := db.GetUser(c, map[string]interface{}{"account": req.Account, "password": req.Password})
	if err != nil {
		logs.CtxWarn(c, "failed to get user, err=%v", err)
		OnFailWithMessage(c, constant.FailedToProcess, "Failed to check existence")
		return
	}
	if len(users) == 0 {
		logs.CtxWarn(c, "length of users is zero")
		OnFailWithMessage(c, constant.FailedToProcess, "The password is wrong/ The user does not exist")
		return
	}

	// Set token
	service.SetToken(c, req.Account, req.RememberMe, c.GetHeader(constant.RealIP))
	OnSuccess(c, dto.SessionManageLoginResponse{NickName: users[0].NickName})
}

func SessionManageSignup(c *gin.Context) {
	req := dto.SessionManageSignupRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))
	// Parameter check
	if req.Account == "" || req.Password == "" || req.NickName == "" || req.Gender == 0 {
		logs.CtxWarn(c, "required parameters are missing")
		OnFail(c, constant.RequestParameterError)
		return
	}

	// Existence check
	if exist, err := service.CheckUserExist(c, req.Account); err != nil {
		logs.CtxWarn(c, "failed to check if user exist, err=%v", err)
		OnFail(c, constant.ServiceError)
		return
	} else if exist {
		logs.CtxWarn(c, "user has existed")
		OnFailWithMessage(c, constant.RequestParameterError, "user has existed")
		return
	}

	// Add to DB
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

func SessionManageLogout(c *gin.Context) {
	req := dto.SessionManageLogoutRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))

	// Parameter check
	if req.Account == "" {
		logs.CtxWarn(c, "required parameters are missing")
		OnFail(c, constant.RequestParameterError)
		return
	}

	// Remove Cookie
	c.SetCookie(constant.Token, "", -1, "/", "", false, false)

	OnSuccess(c, nil)
}
