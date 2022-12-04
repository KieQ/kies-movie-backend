package handler

import (
	"github.com/gin-gonic/gin"
	"kies-movie-backend/constant"
	"kies-movie-backend/i18n"
	"net/http"
)

type Response struct {
	StatusCode    int32       `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data"`
}

func OnSuccess(c *gin.Context, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	successResponse := Response{Data: data}
	c.JSON(http.StatusOK, successResponse)
}

func OnFail(c *gin.Context, statusCode constant.StatusCode) {
	failResponse := Response{StatusCode: int32(statusCode), StatusMessage: i18n.Translate(c, i18n.SentenceIndex(statusCode))}
	c.JSON(http.StatusOK, failResponse)
}

func OnFailWithMessage(c *gin.Context, statusCode constant.StatusCode, index i18n.SentenceIndex) {
	failResponse := Response{StatusCode: int32(statusCode), StatusMessage: i18n.Translate(c, index)}
	c.JSON(http.StatusOK, failResponse)
}
