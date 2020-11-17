package model

import "gorm.io/gorm"

// User 用户类
type User struct {
	gorm.Model
	Username   string      `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password   string      `gorm:"type:varchar(255);not null"`
	Sex        uint8       `gorm:"type:tinyint;"`
	Email      string      `gorm:"type:varchar(100);uniqueIndex"`
	Role       string      `gorm:"type:varchar(100)"`
	Strategies []*Strategy `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Stocks     []*Stock    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Orders     []*Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
