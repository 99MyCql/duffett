package strategy

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/99MyCql/duffett/pkg"
)

func setup() {
	os.Chdir("../../")
	pkg.InitConfig("conf.yaml")
	pkg.InitLog(pkg.DebugLevel)
	pkg.InitDB()
}

func Test_ExecStrategy(t *testing.T) {
	setup()
	filepath := ""
	strategyRsp := ExecStrategy(&filepath, "admin_简单的买卖策略2", "000001.SZ")
	t.Log(strategyRsp)
	log.Debug(os.Remove(filepath))
}
