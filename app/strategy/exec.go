package strategy

// ExecStrategy 执行决策
func ExecStrategy(strategyName string, tsCode string) (float64, error) {
	var (
		amount int
		err    error
	)
	switch strategyName {
	case "随机策略":
		amount, err = randStrategy(tsCode)
	}
	return float64(amount), err
}
