package model

import "gorm.io/gorm"

// Order 订单类
type Order struct {
	gorm.Model
	Money   string `gorm:"type:varchar(255);not null"`
	UserID  uint
	StockId uint
}
