package model

import (
	log "github.com/sirupsen/logrus"
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

func Update(strategy *Strategy) pkg.RspData {
	return pkg.ComUpdate(strategy)
}

func FindById(id uint) *Strategy {
	var s Strategy
	result := pkg.DB.Where("id = ?", id).Find(&s)
	if result.Error != nil || result.RowsAffected < 1 {
		log.Error(result.Error)
		return nil
	}
	return &s
}

func FindByName(strategyName string) *Strategy {
	var s Strategy
	result := pkg.DB.Where("name = ?", strategyName).Find(&s)
	if result.Error != nil || result.RowsAffected < 1 {
		log.Error(result.Error)
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

// UnscopedDeleteById 永久删除
func UnscopedDeleteById(strategyId uint) pkg.RspData {
	result := pkg.DB.Unscoped().Where("id = ?", strategyId).Delete(&Strategy{})
	if result.Error != nil {
		log.Error(result.Error.Error())
		return pkg.ServerErr("服务端删除数据时发生了一些错误")
	}
	return pkg.Suc("")
}
