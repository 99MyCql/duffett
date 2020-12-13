package model

import (
	"gorm.io/gorm"

	"duffett/app/order/model"
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
	orders      []*model.Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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

func FindMonitoringStocks(userID uint) []*Stock {
	var stocks []*Stock = make([]*Stock, 0)
	result := pkg.DB.Where("user_id = ? and state = \"监听中\"", userID).Find(&stocks)
	if result.RowsAffected < 1 {
		return stocks
	}
	return stocks
}
