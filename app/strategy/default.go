package strategy

import (
	"math/rand"
	"time"
)

func randStrategy(tsCode string) (int, error) {
	rand.Seed(time.Now().Unix())
	return rand.Intn(100) - 50, nil
}
