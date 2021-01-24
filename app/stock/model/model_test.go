package model

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"duffett/pkg"
)

func setup() {
	os.Chdir("../../")
	pkg.InitConfig("conf.yaml")
	pkg.InitLog(pkg.DebugLevel)
	pkg.InitDB()
}

func Test_ListMonitoringStocks(t *testing.T) {
	setup()
	stocks := ListMonitoringStocks("admin")
	for i := 0; i < len(stocks); i++ {
		log.Debug(stocks[i])
	}
}
