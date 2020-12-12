package monitor

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"duffett/app/data"
	"duffett/app/order"
	"duffett/app/stock"
	"duffett/app/strategy"
	"duffett/app/trade"
	"duffett/app/user"
	"duffett/pkg"
)

// monitor 监听器
type monitor struct {
	username     string
	tsCode       string
	stockName    string
	strategyName string
	monitorFreq  int64
	ticker       *time.Ticker
	stopped      bool
	mutex        sync.Mutex
	ws           *websocket.Conn
}

// newMonitor 创建一个监听器
func newMonitor(username string, tsCode string, strategyName string, freq int64, ws *websocket.Conn) *monitor {
	if s := strategy.FindByName(strategyName); s == nil {
		ws.WriteJSON(pkg.ClientErr("策略名不存在"))
		return nil
	}
	stockName, err := data.GetStockName(tsCode)
	if err != nil {
		ws.WriteJSON(pkg.ClientErr("tsCode 错误"))
		return nil
	}
	return &monitor{
		username:     username,
		tsCode:       tsCode,
		stockName:    stockName,
		strategyName: strategyName,
		monitorFreq:  freq,
		ticker:       time.NewTicker(time.Second * time.Duration(freq)),
		stopped:      false,
		mutex:        sync.Mutex{},
		ws:           ws,
	}
}

// orderPro 前端所需的 order 数据
type orderPro struct {
	Money        float64
	Price        float64
	State        string
	StockName    string
	StrategyName string
	CreatedAt    string
	UpdatedAt    string
}

// monitoring 启动监听器
func (m *monitor) monitoring() {
	// 在 stock 数据表中记录，监听器结束时删除
	u := user.FindByName(m.username)
	sto := stock.Stock{
		TsCode:      m.tsCode,
		Name:        m.stockName,
		State:       stock.MonitoringState,
		MonitorFreq: m.monitorFreq,
		Share:       0,
		SumProfit:   0,
		CurProfit:   0,
		UserID:      u.ID,
		StrategyID:  strategy.FindByName(m.strategyName).ID,
	}
	stock.Create(&sto)
	defer stock.Delete(&sto)

	// 启动监听器
	for {
		// 判断监听器是否停止
		m.mutex.Lock()
		if m.stopped == true {
			m.mutex.Unlock()
			break
		}
		m.mutex.Unlock()

		// 决策
		amount, err := strategy.ExecStrategy(m.strategyName, m.tsCode)
		if err != nil {
			log.Print(err)
			if err := m.ws.WriteJSON(pkg.ServerErr("服务端决策出错")); err != nil {
				log.Print(err)
			}
			break
		}
		o := orderPro{
			Money:        amount,
			Price:        0,
			State:        order.TradingState,
			StockName:    sto.Name,
			StrategyName: m.strategyName,
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		}
		log.Print(o)
		m.ws.WriteJSON(pkg.SucWithData("", o))

		// 交易
		if err := trade.ExecTrade(m.tsCode, amount); err != nil {
			log.Print(err)
			o.State = order.ErrorState
			m.ws.WriteJSON(pkg.SucWithData("", o))
			break
		}
		o.State = order.TradedState
		o.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		log.Print(o)
		m.ws.WriteJSON(pkg.SucWithData("", o))
		order.Create(&order.Order{
			Money:   o.Money,
			Price:   o.Price,
			State:   o.State,
			UserID:  u.ID,
			StockId: sto.ID,
		})

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
