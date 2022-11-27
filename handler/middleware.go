package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-movie-backend/constant"
	"kies-movie-backend/service"
	"reflect"
	"time"
)

func MiddlewareMetaInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		if logID := c.GetHeader(constant.LogID); logID != "" {
			c.Set(constant.LogID, logID)
		}
	}
}

func MiddlewareAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(constant.Token)
		if err != nil {
			logs.CtxWarn(c, "failed to get token, err=%v", err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		}
		claims, err := service.ValidateToken(tokenStr)
		if err != nil {
			logs.CtxWarn(c, "failed to validate token, err=%v", err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		}
		if val, exist := claims[constant.Account]; !exist {
			logs.CtxWarn(c, "JWT does not contain account")
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		} else {
			c.Set(constant.Account, val)
		}

		if rmAny, exist := claims[constant.RememberMe]; !exist {
			logs.CtxWarn(c, "%v should exist", constant.RememberMe)
			OnFail(c, constant.ServiceError)
			c.Abort()
			return
		} else if rm, ok := rmAny.(bool); !ok {
			logs.CtxWarn(c, "%v should be of type boolean", constant.RememberMe)
			OnFail(c, constant.ServiceError)
			c.Abort()
			return
		} else if rm {
			if expAny, exist := claims["exp"]; !exist {
				logs.CtxWarn(c, "exp should be in JWT when %v is true", constant.RememberMe)
				OnFail(c, constant.ServiceError)
				c.Abort()
				return
			} else if exp, ok := expAny.(float64); !ok {
				logs.CtxWarn(c, "the value of exp is not of type int64, it's %v", reflect.TypeOf(expAny))
				OnFail(c, constant.ServiceError)
				c.Abort()
				return
			} else {
				now := time.Now().Unix()
				if now < int64(exp) && int64(exp)-now < int64(constant.RefreshLimit.Seconds()) {
					service.SetToken(c, c.GetString(constant.Account), rm)
				}
			}
		}
	}
}
