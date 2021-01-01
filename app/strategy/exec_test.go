package strategy

import (
	"os"
	"testing"

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
	if strategyRsp.Code == 0 {
		os.Remove(strategyRsp.Data.(map[string]interface{})["filepath"].(string))
	}
}
