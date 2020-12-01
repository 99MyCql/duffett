package app

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"duffett/pkg"
)

// monitor 监听器
type monitor struct {
	ticker *time.Ticker
	stop   bool
	mutex  sync.Mutex
	ws     *websocket.Conn
}

// newMonitor 创建一个监听器
func newMonitor(freq int64, ws *websocket.Conn) *monitor {
	return &monitor{
		ticker: time.NewTicker(time.Second * time.Duration(freq)),
		stop:   false,
		mutex:  sync.Mutex{},
		ws:     ws,
	}
}

// startMonitor 启动一个监听器
func startMonitor(username string, tsCode string, strategyName string) {
	for {
		// 判断监听器是否停止
		monitors[username][tsCode].mutex.Lock()
		if monitors[username][tsCode].stop == true {
			monitors[username][tsCode].mutex.Unlock()
			monitors[username][tsCode].ws.Close()
			delete(monitors[username], tsCode)
			break
		}
		monitors[username][tsCode].mutex.Unlock()

		// 决策
		amount, err := execStrategy(strategyName, tsCode)
		if err != nil {
			log.Print(err)
			if err := monitors[username][tsCode].ws.WriteJSON(pkg.ServerErr("服务端决策出错")); err != nil {
				log.Print(err)
			}
			break
		}
		log.Print(amount)

		// 交易
		if err := execTrade(tsCode, amount); err != nil {
			log.Print(err)
			if err := monitors[username][tsCode].ws.WriteJSON(pkg.ServerErr("服务端交易出错")); err != nil {
				log.Print(err)
			}
			break
		}

		if err := monitors[username][tsCode].ws.WriteJSON(pkg.SucWithData("交易成功", amount)); err != nil {
			log.Print(err)
		}

		// 等待定时器
		<-monitors[username][tsCode].ticker.C
	}
}

var (
	monitors = make(map[string]map[string]*monitor) // 每个用户都可监视多个股票
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
		log.Print(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Print(req)

	// 获取 username （经过 jwt 中间件时已从 token 中获取）
	username, _ := c.Get("username")

	// 为每个用户创建 map[string]*monitor （如果不存在的话）
	if _, ok := monitors[username.(string)]; !ok {
		monitors[username.(string)] = make(map[string]*monitor)
	}

	// 切换为 websocket 连接
	ws, err := pkg.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusOK, pkg.ServerErr("升级websocket失败"))
		return
	}

	// 创建一个监听器
	monitors[username.(string)][req.TsCode] = newMonitor(req.MonitorFreq, ws)

	// 启动一个监听器
	go startMonitor(username.(string), req.TsCode, req.StrategyName)
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
		log.Print(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Print(req)

	username, _ := c.Get("username")

	// 结束定时器
	monitors[username.(string)][req.TsCode].mutex.Lock()
	monitors[username.(string)][req.TsCode].stop = true
	monitors[username.(string)][req.TsCode].mutex.Unlock()

	c.JSON(http.StatusOK, pkg.Suc("停止监听成功"))
}
