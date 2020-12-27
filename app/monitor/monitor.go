package monitor

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"duffett/app/data"
	orderModel "duffett/app/order/model"
	stockModel "duffett/app/stock/model"
	"duffett/app/strategy"
	strategyModel "duffett/app/strategy/model"
	"duffett/app/trade"
	userModel "duffett/app/user/model"
	"duffett/pkg"
)

// monitor 监听器
type monitor struct {
	// 相关数据信息
	username     string
	strategyName string
	stock        stockModel.Stock
	// 监听所需信息
	ticker  *time.Ticker
	stopped bool
	mutex   sync.Mutex
	ws      *websocket.Conn
}

// newMonitor 创建一个监听器
func newMonitor(username string, tsCode string, strategyName string, freq int64, ws *websocket.Conn) *monitor {
	u := userModel.FindByName(username)
	if u == nil {
		ws.WriteJSON(pkg.ServerErr("查找用户时出错"))
		return nil
	}
	s := strategyModel.FindByName(strategyName)
	if s == nil {
		ws.WriteJSON(pkg.ClientErr("策略名不存在"))
		return nil
	}
	stockName, err := data.GetStockName(tsCode)
	if err != nil {
		log.Print(err)
		ws.WriteJSON(pkg.ClientErr("tsCode 错误"))
		return nil
	}
	return &monitor{
		username:     username,
		strategyName: strategyName,
		stock: stockModel.Stock{
			TsCode:      tsCode,
			Name:        stockName,
			State:       "",
			MonitorFreq: freq,
			Share:       0,
			Profit:      0,
			UserID:      u.ID,
			StrategyID:  s.ID,
		},
		ticker:  time.NewTicker(time.Second * time.Duration(freq)),
		stopped: false,
		mutex:   sync.Mutex{},
		ws:      ws,
	}
}

// orderRsp 前端所需的整合的 order 数据
type orderRsp struct {
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
	m.stock.State = stockModel.MonitoringState
	stockModel.Create(&m.stock)

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
		amount, err := strategy.ExecStrategy(m.strategyName, m.stock.TsCode)
		if err != nil {
			log.Print(err)
			m.ws.WriteJSON(pkg.ServerErr("服务端决策出错"))
			break
		}
		o := orderRsp{
			Money:        amount,
			Price:        0,
			State:        orderModel.TradingState,
			StockName:    m.stock.Name,
			StrategyName: m.strategyName,
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		}
		log.Print(o)
		m.ws.WriteJSON(pkg.SucWithData("", o))

		// 交易
		if err := trade.ExecTrade(m.stock.TsCode, amount); err != nil {
			log.Print(err)
			o.State = orderModel.ErrorState
			o.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			m.ws.WriteJSON(pkg.SucWithData("", o))
			continue
		}
		o.State = orderModel.TradedState
		o.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		log.Print(o)
		m.ws.WriteJSON(pkg.SucWithData("", o))
		orderModel.Create(&orderModel.Order{
			Money:   o.Money,
			Price:   o.Price,
			State:   o.State,
			UserID:  m.stock.UserID,
			StockID: m.stock.ID,
		})

		// 等待定时器
		<-m.ticker.C
	}
}

// finish 结束监听器
func (m *monitor) finish() {
	m.mutex.Lock()
	m.stopped = true
	m.mutex.Unlock()
	// 更新数据库记录
	m.stock.State = stockModel.MonitorFinishState
	stockModel.Update(&m.stock)
}
