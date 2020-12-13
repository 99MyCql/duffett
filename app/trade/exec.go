package trade

import (
	"errors"
	"math/rand"
	"time"
)

// ExecTrade 模拟交易过程
func ExecTrade(TsCode string, amount float64) error {
	time.Sleep(time.Second * 5)
	rand.Seed(time.Now().Unix())
	if rand.Intn(100) > 90 {
		return errors.New("rand error")
	}
	return nil
}
