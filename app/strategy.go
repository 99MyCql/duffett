package app

import (
	"math/rand"
	"time"
)

// ExecStrategy 执行决策
func ExecStrategy(strategyName string, tsCode string) (amount int, err error) {
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
