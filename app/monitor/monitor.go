package monitor

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"duffett/app/data"
	model4 "duffett/app/order/model"
	"duffett/app/stock/model"
	"duffett/app/strategy"
	model2 "duffett/app/strategy/model"
	"duffett/app/trade"
	model3 "duffett/app/user/model"
	"duffett/pkg"
)

// monitor 监听器
type monitor struct {
	// 相关数据信息
	userID       uint
	username     string
	tsCode       string
	stockName    string
	strategyID   uint
	strategyName string
	// 监听所需信息
	monitorFreq int64
	ticker      *time.Ticker
	stopped     bool
	mutex       sync.Mutex
	ws          *websocket.Conn
}

// newMonitor 创建一个监听器
func newMonitor(username string, tsCode string, strategyName string, freq int64, ws *websocket.Conn) *monitor {
	u := model3.FindByName(username)
	if u == nil {
		ws.WriteJSON(pkg.ServerErr("查找用户时出错"))
		return nil
	}
	s := model2.FindByName(strategyName)
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
		userID:       u.ID,
		username:     u.Username,
		tsCode:       tsCode,
		stockName:    stockName,
		strategyID:   s.ID,
		strategyName: s.Name,
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
	sto := model.Stock{
		TsCode:      m.tsCode,
		Name:        m.stockName,
		State:       model.MonitoringState,
		MonitorFreq: m.monitorFreq,
		Share:       0,
		SumProfit:   0,
		CurProfit:   0,
		UserID:      m.userID,
		StrategyID:  m.strategyID,
	}
	model.Create(&sto)
	defer func() {
		sto.State = model.MonitorFinishState
		model.Delete(&sto)
	}()

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
			m.ws.WriteJSON(pkg.ServerErr("服务端决策出错"))
			break
		}
		o := orderPro{
			Money:        amount,
			Price:        0,
			State:        model4.TradingState,
			StockName:    m.stockName,
			StrategyName: m.strategyName,
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		}
		log.Print(o)
		m.ws.WriteJSON(pkg.SucWithData("", o))

		// 交易
		if err := trade.ExecTrade(m.tsCode, amount); err != nil {
			log.Print(err)
			o.State = model4.ErrorState
			o.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			m.ws.WriteJSON(pkg.SucWithData("", o))
			continue
		}
		o.State = model4.TradedState
		o.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		log.Print(o)
		m.ws.WriteJSON(pkg.SucWithData("", o))
		model4.Create(&model4.Order{
			Money:   o.Money,
			Price:   o.Price,
			State:   o.State,
			UserID:  m.userID,
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
