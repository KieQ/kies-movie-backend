package main

import (
	"fmt"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-movie-backend/handler"
	"kies-movie-backend/handler/middleware"
	"kies-movie-backend/model/db"
	"os"
)

func main() {
	db.MustInit()
	StartServer()
}

func StartServer() {
	var port = os.Getenv("PORT")
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
	g.Use(gin.Logger(), gin.Recovery(), middleware.MetaInfo())

	g.GET("/ping", handler.Ping)

	user := g.Group("/user")
	user.POST("/add", handler.UserAdd)
	user.POST("/update", handler.UserUpdate)
	user.GET("/detail", handler.UserDetail)
	user.GET("/list", handler.UserList)
}
