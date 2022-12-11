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

	sessionManage := g.Group("/session_manage")
	sessionManage.POST("/log_in", handler.SessionManageLogin)
	sessionManage.POST("/sign_up", handler.SessionManageSignup)
	sessionManage.POST("/log_out", handler.MiddlewareAuthority(), handler.SessionManageLogout)

	user := g.Group("/user")
	user.Use(handler.MiddlewareAuthority())
	user.POST("/update", handler.UserUpdate)
	user.GET("/detail", handler.UserDetail)
	user.GET("/list", handler.UserList)

	homepage := g.Group("/homepage")
	homepage.GET("/content", handler.HomepageContent)

	movie := g.Group("/movie")
	movie.Use(handler.MiddlewareAuthority())
	movie.GET("/list", handler.MovieList)
	movie.GET("/detail", handler.MovieDetail)
	movie.POST("/like", handler.MovieLike)
	movie.POST("/update", handler.MovieUpdate)
	movie.POST("/add", handler.MovieAdd)
	movie.POST("/delete", handler.MovieDelete)

	movieNotLogin := g.Group("/movie/not_login")
	movieNotLogin.GET("/list", handler.NotLoginMovieList)
	movieNotLogin.GET("/detail", handler.NotLoginMovieDetail)

}
