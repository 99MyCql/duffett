package main

import (
	log "github.com/sirupsen/logrus"

	orderModel "duffett/app/order/model"
	stockModel "duffett/app/stock/model"
	strategyModel "duffett/app/strategy/model"
	userModel "duffett/app/user/model"
	"duffett/pkg"
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
