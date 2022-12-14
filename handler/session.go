package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-movie-backend/constant"
	"kies-movie-backend/dto"
	"kies-movie-backend/i18n"
	"kies-movie-backend/model/db"
	"kies-movie-backend/model/table"
	"kies-movie-backend/service"
	"kies-movie-backend/utils"
	"time"
)

func SessionLogin(c *gin.Context) {
	req := dto.SessionLoginRequest{}
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
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToCheckExistence)
		return
	}
	if len(users) == 0 {
		logs.CtxWarn(c, "length of users is zero")
		OnFailWithMessage(c, constant.FailedToProcess, i18n.FailedToLogin)
		return
	}

	// Set token
	service.SetToken(c, req.Account, req.RememberMe, c.GetHeader(constant.RealIP))
	OnSuccess(c, dto.SessionLoginResponse{NickName: users[0].NickName})
}

func SessionSignup(c *gin.Context) {
	req := dto.SessionSignupRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logs.CtxWarn(c, "failed to bind request, err=%v", err)
		OnFail(c, constant.RequestParameterError)
		return
	}
	logs.CtxInfo(c, "req=%v", utils.ToJSON(req))
	// Parameter check
	if req.Account == "" || req.Password == "" || req.NickName == "" {
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
		OnFailWithMessage(c, constant.RequestParameterError, i18n.UserHasExisted)
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
		DefaultLanguage:  req.DefaultLanguage,
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

func SessionLogout(c *gin.Context) {
	req := dto.SessionLogoutRequest{}
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
	if c.GetString(constant.Account) != req.Account {
		logs.CtxWarn(c, "authority check failed")
		OnFail(c, constant.NoAuthority)
		return
	}

	// Remove Cookie
	c.SetCookie(constant.Token, "", -1, "/", "", false, false)

	OnSuccess(c, nil)
}
