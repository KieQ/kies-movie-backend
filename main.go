package main

import (
	"fmt"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-movie-backend/handler"
	"kies-movie-backend/model/db"
	"kies-movie-backend/utils"
	"os"
)

func main() {
	utils.InitRandom()
	db.MustInit()
	StartServer()
}

func StartServer() {
	var port = os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	logs.Info("Server start with port %v", port)
	server := gin.New()
	Register(server)
	if err := server.Run(fmt.Sprintf(":%v", port)); err != nil {
		panic(err)
	}
}

func Register(g *gin.Engine) {
	g.Use(gin.Logger(), gin.Recovery(), handler.MiddlewareMetaInfo())

	g.GET("/ping", handler.Ping)

	session := g.Group("/session")
	session.POST("/log_in", handler.SessionLogin)
	session.POST("/sign_up", handler.SessionSignup)
	session.POST("/log_out", handler.MiddlewareAuthority(), handler.SessionLogout)

	user := g.Group("/user")
	user.Use(handler.MiddlewareAuthority())
	user.POST("/update", handler.UserUpdate)
	user.GET("/detail", handler.UserDetail)
	user.GET("/list", handler.UserList)

	homepage := g.Group("/homepage")
	homepage.GET("/content", handler.HomepageContent)

	video := g.Group("/video")
	video.Use(handler.MiddlewareAuthority())
	video.GET("/list", handler.VideoList)
	video.GET("/detail", handler.VideoDetail)
	video.POST("/like", handler.VideoLike)
	video.POST("/update", handler.VideoUpdate)
	video.POST("/add", handler.VideoAdd)
	video.POST("/delete", handler.VideoDelete)

	videoNotLogin := g.Group("/video/not_login")
	videoNotLogin.GET("/list", handler.NotLoginVideoList)
	videoNotLogin.GET("/detail", handler.NotLoginVideoDetail)

}
