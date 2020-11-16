package model

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255);not null;uniqueIndex"`
	State     string `gorm:"type:varchar(100);mot null"`
	SumProfit string `gorm:"type:varchar(255)"`
	CurProfit string `gorm:"type:varchar(255)"`
	UserID    uint
	orders    []*Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
