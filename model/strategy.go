package model

import "gorm.io/gorm"

type Strategy struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Desc    string `gorm:"type:varchar(255)"`
	Content string `gorm:"type:text"`
	UserID  uint
}
