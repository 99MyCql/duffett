package app

import (
	"math/rand"
	"time"
)

// execStrategy 执行决策
func execStrategy(strategyName string, tsCode string) (amount int, err error) {
	switch strategyName {
	default:
		amount, err = randStrategy(tsCode)
	}
	return amount, err
}

func randStrategy(tsCode string) (int, error) {
	rand.Seed(time.Now().Unix())
	return rand.Intn(100) - 50, nil
}
