package utils

import (
	"math/rand"
	"time"
)

func InitRandom() {
	rand.Seed(time.Now().Unix())
}
