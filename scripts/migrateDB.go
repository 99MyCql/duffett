package main

import (
	log "github.com/sirupsen/logrus"

	orderModel "github.com/99MyCql/duffett/app/order/model"
	stockModel "github.com/99MyCql/duffett/app/stock/model"
	strategyModel "github.com/99MyCql/duffett/app/strategy/model"
	userModel "github.com/99MyCql/duffett/app/user/model"
	"github.com/99MyCql/duffett/pkg"
)

func init() {
	pkg.InitConfig("conf.yaml")
	pkg.InitLog(pkg.DebugLevel)
	pkg.InitDB()
}

func main() {
	// 自动创建数据表
	err := pkg.DB.AutoMigrate(&userModel.User{}, &strategyModel.Strategy{}, &stockModel.Stock{}, &orderModel.Order{})
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("migrate successfully")
}
