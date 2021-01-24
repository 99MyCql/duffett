package model

import (
	"gorm.io/gorm"

	"duffett/pkg"
)

// Order 订单类
type Order struct {
	gorm.Model
	Money   float64 `gorm:"type:double;not null"`
	Price   float64 `gorm:"type:double;not null"`
	State   string  `gorm:"type:varchar(100);not null"`
	UserID  uint
	StockID uint
}

const (
	TradingState   = "报单中"
	TradedState    = "已成交"
	CancelingState = "撤单中"
	CancelledState = "已撤单"
	ErrorState     = "出错"
)

func Create(order *Order) pkg.RspData {
	return pkg.ComCreate(order)
}

func Delete(order *Order) pkg.RspData {
	return pkg.ComDelete(order)
}

func Update(order *Order) pkg.RspData {
	return pkg.ComUpdate(order)
}

// ListOrders 与 user 表、stock 表连接查询
func ListOrders(username string, stockID uint) []map[string]interface{} {
	orders := make([]map[string]interface{}, 0)
	pkg.DB.
		Table("orders").
		Select("orders.id, orders.money, orders.price, orders.state, orders.created_at").
		Joins("JOIN users ON users.id = orders.user_id").
		Joins("JOIN stocks ON stocks.id = orders.stock_id").
		Where("users.username = ?", username).
		Where("stocks.id = ?", stockID).
		Scan(&orders)
	return orders
}
