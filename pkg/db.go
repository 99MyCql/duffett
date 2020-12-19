package pkg

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 数据库操作对象
var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() {
	var err error
	// 创建数据库连接池
	DB, err = gorm.Open(mysql.Open(Conf.MysqlUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

// ComCreate 通用数据库创建
func ComCreate(model interface{}) RspData {
	result := DB.Create(model)
	if result.Error != nil {
		log.Print(result.Error.Error())
		return ServerErr("服务端创建数据时发生了一些错误")
	}
	return Suc("")
}

// ComDelete 通用数据库删除
func ComDelete(model interface{}) RspData {
	result := DB.Delete(model)
	if result.Error != nil {
		log.Print(result.Error.Error())
		return ServerErr("服务端删除数据时发生了一些错误")
	}
	return Suc("")
}

// ComUpdate 通用数据库更新
func ComUpdate(model interface{}) RspData {
	result := DB.Save(model)
	if result.Error != nil {
		log.Print(result.Error.Error())
		return ServerErr("服务端更新数据时发生了一些错误")
	}
	return Suc("")
}
