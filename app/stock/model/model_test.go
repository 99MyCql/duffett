package model

import (
	"log"
	"testing"

	"duffett/pkg"
)

func setup() {
	pkg.InitConfig("../../conf.yaml")
	pkg.InitLog()
	pkg.InitDB()
}

func Test_findMonitoringStocks(t *testing.T) {
	setup()
	stocks := findMonitoringStocks(1)
	for i := 0; i < len(stocks); i++ {
		log.Print(stocks[i])
	}
}
