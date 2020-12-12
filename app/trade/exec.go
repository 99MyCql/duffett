package trade

import "time"

func ExecTrade(TsCode string, amount float64) error {
	time.Sleep(time.Second * 10)
	return nil
}
