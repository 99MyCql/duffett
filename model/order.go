package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Money   string `gorm:"type:varchar(255);not null"`
	UserID  uint
	StockId uint
}
