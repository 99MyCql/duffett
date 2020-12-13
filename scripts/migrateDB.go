package main

import (
	"log"

	model4 "duffett/app/order/model"
	"duffett/app/stock/model"
	model2 "duffett/app/strategy/model"
	model3 "duffett/app/user/model"
	"duffett/pkg"
)

func init() {
	pkg.InitConfig("conf.yaml")
	pkg.InitLog()
	pkg.InitDB()
}

func main() {
	// 自动创建数据表
	err := pkg.DB.AutoMigrate(&model3.User{}, &model2.Strategy{}, &model.Stock{}, &model4.Order{})
	if err != nil {
		log.Fatal(err)
	}
}
