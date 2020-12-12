package user

import (
	"errors"
	"log"

	"gorm.io/gorm"

	"duffett/app/order"
	"duffett/app/stock"
	"duffett/app/strategy"
	"duffett/pkg"
)

// User 用户类
type User struct {
	gorm.Model
	Username   string               `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password   string               `gorm:"type:varchar(255);not null"`
	Sex        uint8                `gorm:"type:tinyint;"`
	Email      string               `gorm:"type:varchar(100);uniqueIndex"`
	Role       string               `gorm:"type:varchar(100)"`
	Strategies []*strategy.Strategy `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Stocks     []*stock.Stock       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Orders     []*order.Order       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// check 检查用户名和密码是否正确
func check(username string, password string) pkg.RspData {
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

func create(user *User) pkg.RspData {
	return pkg.ComCreate(user)
}
