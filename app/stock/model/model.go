package model

import (
	"gorm.io/gorm"

	"github.com/99MyCql/duffett/app/order/model"
	"github.com/99MyCql/duffett/pkg"
)

// Stock 监听股票类
type Stock struct {
	gorm.Model
	TsCode      string  `gorm:"type:varchar(255);not null"`
	Name        string  `gorm:"type:varchar(255);not null"`
	State       string  `gorm:"type:varchar(100);not null"`
	MonitorFreq int64   `gorm:"type:bigint;not null"`
	Share       float64 `gorm:"type:double;not null"`
	Profit      float64 `gorm:"type:double"`
	UserID      uint
	StrategyID  uint
	Orders      []*model.Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

const (
	MonitoringState    = "监听中"
	MonitorFinishState = "监听结束"
)

func Create(stock *Stock) pkg.RspData {
	return pkg.ComCreate(stock)
}

func Delete(stock *Stock) pkg.RspData {
	return pkg.ComDelete(stock)
}

func Update(stock *Stock) pkg.RspData {
	return pkg.ComUpdate(stock)
}

// ListMonitoringStocks 与 user 表连接查询监听中的股票
func ListMonitoringStocks(username string) []*Stock {
	stocks := make([]*Stock, 0)
	pkg.DB.
		Where("stocks.state = \"监听中\"").
		Joins("JOIN users ON users.id = stocks.user_id").
		Where("users.username = ?", username).
		Find(&stocks)
	return stocks
}

// ListStocks 与 user 表连接查询所有记录的股票
func ListStocks(username string) []map[string]interface{} {
	stockPros := make([]map[string]interface{}, 0)
	pkg.DB.
		Table("stocks").
		Select("stocks.id, stocks.ts_code, stocks.name, stocks.state, stocks.monitor_freq, stocks.share, "+
			"stocks.profit, stocks.created_at, stocks.updated_at, "+
			"strategies.name as strategyName").
		Joins("JOIN users ON users.id = stocks.user_id").
		Joins("JOIN strategies ON strategies.id = stocks.strategy_id").
		Where("users.username = ?", username).
		Scan(&stockPros)
	return stockPros
}
