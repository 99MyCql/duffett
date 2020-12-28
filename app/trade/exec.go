package trade

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/99MyCql/duffett/app/data"
)

// ExecTrade 模拟交易过程
func ExecTrade(tsCode string, amount float64) (float64, error) {
	time.Sleep(time.Second * 5)
	rand.Seed(time.Now().Unix())
	if rand.Intn(100) > 90 {
		return 0, errors.New("rand error")
	}
	realTimeData, err := data.GetRealTimeData(tsCode)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return realTimeData.CurPrice, nil
}
