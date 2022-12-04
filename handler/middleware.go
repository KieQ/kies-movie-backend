package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-movie-backend/constant"
	"kies-movie-backend/i18n"
	"kies-movie-backend/service"
	"kies-movie-backend/utils"
	"time"
)

func MiddlewareMetaInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.RequestID, c.GetHeader(constant.RequestID))
		c.Set(constant.RealIP, c.GetHeader(constant.RealIP))
		c.Set(i18n.ContextLanguage, c.Query(constant.Lang))
	}
}

func MiddlewareAuthority(allowNotLogin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get Token from cookie
		tokenStr, err := c.Cookie(constant.Token)
		if err != nil {
			if allowNotLogin {
				logs.CtxInfo(c, "this request is not from logged in user but allowed")
				c.Set(constant.NotLogin, true)
				return
			} else {
				logs.CtxWarn(c, "failed to get token, err=%v", err)
				OnFail(c, constant.UserNotLogin)
				c.Abort()
				return
			}
		}

		claims, err := service.ValidateToken(tokenStr)
		if err != nil {
			logs.CtxWarn(c, "failed to validate token, err=%v", err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		}

		if val, err := utils.GetFromAnyMap[string](claims, constant.Account); err != nil {
			logs.CtxWarn(c, "JWT does not contain %v, err=%v", constant.Account, err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		} else {
			c.Set(constant.Account, val)
		}

		if val, err := utils.GetFromAnyMap[string](claims, constant.TokenIP); err != nil {
			logs.CtxWarn(c, "JWT does not contain %v, err=%v", constant.Account, err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		} else if val != c.GetHeader(constant.RealIP) {
			logs.CtxWarn(c, "user ip has changed from %v to %v", val, c.GetHeader(constant.RealIP))
			c.SetCookie(constant.Token, "", -1, "/", "", false, false)
			OnFail(c, constant.UserIPChanged)
			c.Abort()
			return
		}

		if rm, err := utils.GetFromAnyMap[bool](claims, constant.RememberMe); err != nil {
			logs.CtxWarn(c, "failed to get %v, err=%v", constant.RememberMe, err)
			OnFail(c, constant.ServiceError)
			c.Abort()
			return
		} else if rm {
			if exp, err := utils.GetFromAnyMap[float64](claims, "exp"); err != nil {
				logs.CtxWarn(c, "failed to get exp, err=%v", err)
				OnFail(c, constant.ServiceError)
				c.Abort()
				return
			} else {
				now := time.Now().Unix()
				if now < int64(exp) && int64(exp)-now < int64(constant.RefreshLimit.Seconds()) {
					service.SetToken(c, c.GetString(constant.Account), rm, c.GetHeader(constant.RealIP))
				}
			}
		}

	}
}
