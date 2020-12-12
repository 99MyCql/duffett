package stock

import (
	"gorm.io/gorm"

	"duffett/app/order"
	"duffett/pkg"
)

// Stock 股票类
type Stock struct {
	gorm.Model
	TsCode      string  `gorm:"type:varchar(255);not null"`
	Name        string  `gorm:"type:varchar(255);not null"`
	State       string  `gorm:"type:varchar(100);not null"`
	MonitorFreq int64   `gorm:"type:bigint;not null"`
	Share       float64 `gorm:"type:double;not null"`
	SumProfit   float64 `gorm:"type:double"`
	CurProfit   float64 `gorm:"type:double"`
	UserID      uint
	StrategyID  uint
	orders      []*order.Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

const (
	MonitoringState = "监听中"
)

func Create(stock *Stock) pkg.RspData {
	return pkg.ComCreate(stock)
}

func Delete(stock *Stock) pkg.RspData {
	return pkg.ComDelete(stock)
}
