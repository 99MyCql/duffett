package pkg

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"duffett/model"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// 创建数据库连接池
	DB, err = gorm.Open(mysql.Open(Conf.MysqlUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 自动创建数据表
	err = DB.AutoMigrate(&model.User{}, &model.Strategy{}, &model.Stock{}, &model.Order{})
	if err != nil {
		log.Fatalln(err)
	}
}
