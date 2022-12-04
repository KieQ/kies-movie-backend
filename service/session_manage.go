package service

import (
	"context"
	"errors"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"kies-movie-backend/constant"
	"kies-movie-backend/model/db"
	"os"
	"time"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

func CheckUserExist(ctx context.Context, account string) (bool, error) {
	logs.CtxInfo(ctx, "check account for %v", account)

	users, err := db.GetUser(ctx, map[string]interface{}{"account": account})
	if err != nil {
		logs.CtxWarn(ctx, "failed to check, err=%v", err)
		return false, err
	}
	return len(users) > 0, nil
}

func SetToken(c *gin.Context, account string, rememberMe bool, ip string) {
	var maxAge = 0
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	if rememberMe {
		maxAge = int(constant.RememberMeDuration.Seconds())
		claims["exp"] = time.Now().Add(constant.RememberMeDuration).Unix()
	}
	claims[constant.Account] = account
	claims[constant.TokenIP] = ip

	s, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logs.CtxWarn(c, "failed to signed the jwt, err", err)
		return
	}
	c.SetCookie(constant.Token, s, maxAge, "/", "", false, false)
}

func ValidateToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid JWT method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Valid {
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			return mapClaims, nil
		}
		return nil, errors.New("claims is not map")
	} else {
		return nil, errors.New("validation has encountered unreachable code")
	}
}
