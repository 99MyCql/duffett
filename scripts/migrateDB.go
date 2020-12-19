package main

import (
	"log"

	orderModel "duffett/app/order/model"
	stockModel "duffett/app/stock/model"
	strategyModel "duffett/app/strategy/model"
	userModel "duffett/app/user/model"
	"duffett/pkg"
)

func init() {
	pkg.InitConfig("conf.yaml")
	pkg.InitLog()
	pkg.InitDB()
}

func main() {
	// 自动创建数据表
	err := pkg.DB.AutoMigrate(&userModel.User{}, &strategyModel.Strategy{}, &stockModel.Stock{}, &orderModel.Order{})
	if err != nil {
		log.Fatal(err)
	}
	log.Print("migrate successfully")
}
