package app

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"duffett/pkg"
)

type monitor struct {
	ticker *time.Ticker
	stop   bool
	mutex  sync.Mutex
}

var (
	timeTickers = make(map[string]map[string]*monitor)
)

type startMonitorReq struct {
	TsCode       string `json:"ts_code" binding:"required,excludes= "`       // 股票代码
	StrategyName string `json:"strategy_name" binding:"required,excludes= "` // 策略名字
	MonitorFreq  int64  `json:"monitor_freq" binding:"required,number"`      // 监听频率，以秒为单位
}

// @Summary StartMonitor
// @Tags Monitor
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Param data body startMonitorReq true "data"
// @Success 200 {string} json "{"code":x,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":x,"data":{},"msg":""}"
// @Router /api/v1/monitor/start [post]
func StartMonitor(c *gin.Context) {
	var req startMonitorReq
	if err := c.ShouldBind(&req); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}

	// 获取 username （经过 jwt 中间件时已从 token 中获取）
	username, _ := c.Get("username")

	// 启动一个监听器
	go func() {
		timeTickers[username.(string)][req.TsCode] = &monitor{
			ticker: time.NewTicker(time.Second * time.Duration(req.MonitorFreq)), // 创建一个定时器
			stop:   false,
		}

		for {
			// 判断监听器是否停止
			timeTickers[username.(string)][req.TsCode].mutex.Lock()
			if timeTickers[username.(string)][req.TsCode].stop == true {
				timeTickers[username.(string)][req.TsCode].mutex.Unlock()
				break
			}
			timeTickers[username.(string)][req.TsCode].mutex.Unlock()

			// 决策
			tradeCode, amount, err := defaultStrategy1(req.TsCode)
			if err != nil {
				// TODO 决策出错通知前端
				break
			}

			// 交易
			err = execTrade(req.TsCode, tradeCode, amount)
			if err != nil {
				// TODO 交易出错通知前端
				break
			}
			// 等待定时器
			<-timeTickers[username.(string)][req.StrategyName].ticker.C
		}
	}()

	c.JSON(http.StatusOK, pkg.Suc("启动监听成功"))
}

type stopMonitorReq struct {
	TsCode string `json:"ts_code" binding:"required"`
}

// @Summary StopMonitor
// @Tags Monitor
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Param data body stopMonitorReq true "data"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":非0,"data":{},"msg":""}"
// @Router /api/v1/monitor/stop [post]
func StopMonitor(c *gin.Context) {
	var req stopMonitorReq
	if err := c.ShouldBind(&req); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}

	username, _ := c.Get("username")

	// 结束定时器
	timeTickers[username.(string)][req.TsCode].mutex.Lock()
	timeTickers[username.(string)][req.TsCode].stop = true
	timeTickers[username.(string)][req.TsCode].mutex.Unlock()

	c.JSON(http.StatusOK, pkg.Suc("停止监听成功"))
}
