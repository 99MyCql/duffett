package main

import (
	"log"
	"os/exec"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"duffett/app"
	_ "duffett/docs"
	"duffett/middleware"
	"duffett/pkg"
)

func init() {
	pkg.InitConfig()
	pkg.InitLog()
	pkg.InitDB()
}

func main() {
	router := gin.Default()

	// 访问 swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 配置路由
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/login", app.Login)
			user.GET("/test", middleware.JWTAuth(), app.Test)
		}
	}

	// 生成 swagger 文档目录
	cmd := exec.Command("swag", "init")
	out, err := cmd.CombinedOutput()
	log.Print(string(out))
	if err != nil {
		log.Fatal(err)
	}

	// 运行
	router.Run(pkg.Conf.Addr)
}
