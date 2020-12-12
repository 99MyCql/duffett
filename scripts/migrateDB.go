package main

import (
	"log"

	"duffett/app/order"
	"duffett/app/stock"
	"duffett/app/strategy"
	"duffett/app/user"
	"duffett/pkg"
)

func init() {
	pkg.InitConfig("conf.yaml")
	pkg.InitLog()
	pkg.InitDB()
}

func main() {
	// 自动创建数据表
	err := pkg.DB.AutoMigrate(&user.User{}, &strategy.Strategy{}, &stock.Stock{}, &order.Order{})
	if err != nil {
		log.Fatal(err)
	}
}
