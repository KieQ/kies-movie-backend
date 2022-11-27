package db

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
)

func MustInit() {
	if err := InitUserDB(); err != nil {
		logs.Fatal("failed to init UserDB, %v", err)
	}
}
