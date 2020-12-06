package monitor

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"duffett/app"
	"duffett/pkg"
)

// monitor 监听器
type monitor struct {
	username     string
	tsCode       string
	strategyName string
	ticker       *time.Ticker
	stopped      bool
	mutex        sync.Mutex
	ws           *websocket.Conn
}

// newMonitor 创建一个监听器
func newMonitor(username string, tsCode string, strategyName string, freq int64, ws *websocket.Conn) *monitor {
	return &monitor{
		username:     username,
		tsCode:       tsCode,
		strategyName: strategyName,
		ticker:       time.NewTicker(time.Second * time.Duration(freq)),
		stopped:      false,
		mutex:        sync.Mutex{},
		ws:           ws,
	}
}

// start 启动监听器
func (m *monitor) start() {
	for {
		// 判断监听器是否停止
		m.mutex.Lock()
		if m.stopped == true {
			m.mutex.Unlock()
			break
		}
		m.mutex.Unlock()

		// 决策
		amount, err := app.ExecStrategy(m.strategyName, m.tsCode)
		if err != nil {
			log.Print(err)
			if err := m.ws.WriteJSON(pkg.ServerErr("服务端决策出错")); err != nil {
				log.Print(err)
			}
			break
		}
		log.Print(amount)

		// 交易
		if err := app.ExecTrade(m.tsCode, amount); err != nil {
			log.Print(err)
			if err := m.ws.WriteJSON(pkg.ServerErr("服务端交易出错")); err != nil {
				log.Print(err)
			}
			break
		}
		if err := m.ws.WriteJSON(pkg.SucWithData("交易成功", amount)); err != nil {
			log.Print(err)
		}

		// 等待定时器
		<-m.ticker.C
	}
}

// stop 停止监听器
func (m *monitor) stop() {
	m.mutex.Lock()
	m.stopped = true
	m.mutex.Unlock()
}
