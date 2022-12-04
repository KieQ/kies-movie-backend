package handler

import (
	"github.com/gin-gonic/gin"
	"kies-movie-backend/constant"
)

func MovieAll(c *gin.Context) {
	if c.GetBool(constant.NotLogin) {
		MovieAllNotLogin(c)
	} else {
		MovieAllLogin(c)
	}
}

func MovieAllNotLogin(c *gin.Context) {
	OnSuccess(c, nil)
}

func MovieAllLogin(c *gin.Context) {
	OnSuccess(c, nil)
}
