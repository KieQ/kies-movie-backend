package main

import (
	"fmt"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-movie-backend/handler"
	"kies-movie-backend/model/db"
	"os"
)

func main() {
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

	sessionManage := g.Group("/session_manage")
	sessionManage.POST("/log_in", handler.SessionManageLogin)
	sessionManage.POST("/sign_up", handler.SessionManageSignup)
	sessionManage.POST("/log_out", handler.MiddlewareAuthority(false), handler.SessionManageLogout)

	user := g.Group("/user")
	user.Use(handler.MiddlewareAuthority(false))
	user.POST("/update", handler.UserUpdate)
	user.GET("/detail", handler.UserDetail)
	user.GET("/list", handler.UserList)

	movie := g.Group("/movie")
	movie.Use(handler.MiddlewareAuthority(true))
	movie.GET("/list", handler.MovieAll)
}
