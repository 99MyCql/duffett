package strategy

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"duffett/app/strategy/model"
	"duffett/pkg"
)

const (
	codeHeader string = `package main

import (
	"fmt"
	"math/rand"
	"time"

	"duffett/app/data"
	"duffett/pkg"
)

func init() {
	pkg.InitConfig("conf.yaml")
	pkg.InitLog(pkg.FatalLevel)
}
`
)

// ExecStrategy 执行决策
func ExecStrategy(filepath *string, strategyName string, tsCode string) pkg.RspData {
	// 获取决策
	s := model.FindByName(strategyName)
	if s == nil {
		log.Error("错误的决策名：" + strategyName)
		return pkg.ClientErr("错误的决策名")
	}

	if *filepath == "" {
		// 设置文件路径
		*filepath = "temp/" + strategyName + fmt.Sprintf("%v", time.Now().UnixNano()) + ".go"

		// 创建go文件并写入代码
		f, err := os.OpenFile(*filepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Error(err)
			return pkg.ServerErr("")
		}
		defer f.Close()
		_, err = f.WriteString(codeHeader + s.Content)
		if err != nil {
			log.Error(err)
			return pkg.ServerErr("")
		}

	}

	// 设置执行时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 格式化导入包
	cmd := exec.CommandContext(ctx, "goimports", "-w", *filepath)
	out, _ := cmd.CombinedOutput()
	log.Debug(string(out))

	// 运行代码
	cmd = exec.CommandContext(ctx, "go", "run", *filepath)
	stdin, _ := cmd.StdinPipe()
	stdin.Write([]byte(fmt.Sprintln(tsCode))) // 设置输入
	out, err := cmd.CombinedOutput()          // 获取输出
	log.Debug(string(out))
	if err != nil {
		log.Error(err)
		return pkg.ClientErr(string(out))
	}

	amount, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		log.Error(err)
		return pkg.ClientErr(string(out))
	}

	return pkg.SucWithData("", amount)
}

func execStrategy(strategyName string, tsCode string) (float64, error) {
	code := `func main() {
	var tsCode string
	var ans float64

	fmt.Scanf("%s", &tsCode)

	realTime, err := data.GetRealTimeData(tsCode)
	if err != nil {
		fmt.Print(err)
		return
	}

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	yesterdayData, err := data.GetDailyData(tsCode, yesterday)
	if err != nil {
		fmt.Print(err)
		return
	}

	if (realTime.CurPrice-yesterdayData.Close)/yesterdayData.Close >= 0.05 {
		ans = -200
	} else if (realTime.CurPrice-yesterdayData.Close)/yesterdayData.Close <= -0.05 {
		ans = 200
	} else {
		ans = 0
	}
	fmt.Print(ans)
}
`
	code = codeHeader + code

	filepath := "temp/" + strategyName + fmt.Sprintf("%v", time.Now().UnixNano()) + ".go"

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	_, err = f.WriteString(code)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "goimports", "-w", filepath)
	out, _ := cmd.CombinedOutput()
	log.Debug(string(out))

	cmd = exec.CommandContext(ctx, "go", "run", filepath)
	stdin, _ := cmd.StdinPipe()
	stdin.Write([]byte(fmt.Sprintln(tsCode)))
	out, err = cmd.CombinedOutput()
	log.Debug(string(out))
	if err != nil {
		log.Error(err)
		return 0, err
	}
	amount, _ := strconv.ParseFloat(string(out), 64)

	return amount, nil
}
