package model

import (
	"gorm.io/gorm"

	"duffett/app/stock/model"
	"duffett/pkg"
)

// Strategy 策略类
type Strategy struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Desc    string `gorm:"type:varchar(255)"`
	Content string `gorm:"type:text"`
	UserID  uint
	Stocks  []*model.Stock `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func FindByName(strategyName string) *Strategy {
	var s Strategy
	result := pkg.DB.Where("name = ?", strategyName).Find(&s)
	if result.RowsAffected < 1 {
		return nil
	}
	return &s
}
