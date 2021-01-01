package model

import (
	"gorm.io/gorm"

	"github.com/99MyCql/duffett/app/stock/model"
	"github.com/99MyCql/duffett/pkg"
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

func Create(strategy *Strategy) pkg.RspData {
	return pkg.ComCreate(strategy)
}

func Delete(strategy *Strategy) pkg.RspData {
	return pkg.ComDelete(strategy)
}

func Update(strategy *Strategy) pkg.RspData {
	return pkg.ComUpdate(strategy)
}

func FindById(id uint) *Strategy {
	var s Strategy
	result := pkg.DB.Where("id = ?", id).Find(&s)
	if result.RowsAffected < 1 {
		return nil
	}
	return &s
}

func FindByName(strategyName string) *Strategy {
	var s Strategy
	result := pkg.DB.Where("name = ?", strategyName).Find(&s)
	if result.RowsAffected < 1 {
		return nil
	}
	return &s
}

func ListStrategies(username string) []*Strategy {
	strategies := make([]*Strategy, 0)
	pkg.DB.
		Joins("JOIN users ON users.id = strategies.user_id").
		Where("users.username = ?", username).
		Find(&strategies)
	return strategies
}
