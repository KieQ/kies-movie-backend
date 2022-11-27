package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	logs.CtxInfo(c, "ping has been sent")
	OnSuccess(c, nil)
}
