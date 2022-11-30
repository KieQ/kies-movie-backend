package db

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var movieDB *gorm.DB

func MustInit() {
	if err := InitUserDB(); err != nil {
		logs.Fatal("failed to init UserDB, %v", err)
	}
}

func InitUserDB() error {
	dsn := os.Getenv("MOVIE_DB_POSTGRES_URL")
	var err error
	movieDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}
