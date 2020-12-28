package model

import (
	"errors"
	"log"

	"gorm.io/gorm"

	model3 "github.com/99MyCql/duffett/app/order/model"
	"github.com/99MyCql/duffett/app/stock/model"
	model2 "github.com/99MyCql/duffett/app/strategy/model"
	"github.com/99MyCql/duffett/pkg"
)

// User 用户类
type User struct {
	gorm.Model
	Username   string             `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password   string             `gorm:"type:varchar(255);not null"`
	Sex        uint8              `gorm:"type:tinyint;"`
	Email      string             `gorm:"type:varchar(100);uniqueIndex"`
	Role       string             `gorm:"type:varchar(100)"`
	Strategies []*model2.Strategy `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Stocks     []*model.Stock     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Orders     []*model3.Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Check 检查用户名和密码是否正确
func Check(username string, password string) pkg.RspData {
	result := pkg.DB.Where("username = ? and password = ?",
		username, pkg.Md5Encode(password)).First(&User{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return pkg.ClientErr("用户名或密码错误")
	} else if result.Error != nil {
		log.Print(result.Error)
		return pkg.ServerErr("服务端发生了一些错误")
	}
	return pkg.Suc("")
}

func FindByName(username string) *User {
	var user User
	result := pkg.DB.Where("username = ?", username).Find(&user)
	if result.RowsAffected < 1 {
		return nil
	}
	return &user
}

func Create(user *User) pkg.RspData {
	return pkg.ComCreate(user)
}
