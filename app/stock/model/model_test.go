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
	stocks := FindMonitoringStocks("admin")
	for i := 0; i < len(stocks); i++ {
		log.Print(stocks[i])
	}
}
