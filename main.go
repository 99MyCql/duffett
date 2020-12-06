package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"duffett/app"
	appMonitor "duffett/app/monitor"
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
			user.POST("/register", app.Register)
			user.GET("/testJwt", middleware.JWTAuth(), app.TestJwt)
		}

		data := v1.Group("/data").Use(middleware.JWTAuth())
		{
			data.POST("/tushare", app.Tushare)
		}

		monitor := v1.Group("/monitor")
		{
			monitor.GET("/ws", appMonitor.WS)
		}
	}

	// 运行
	router.Run(pkg.Conf.Addr)
}
