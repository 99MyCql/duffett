package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	appData "duffett/app/data"
	appMonitor "duffett/app/monitor"
	appOrder "duffett/app/order"
	appStock "duffett/app/stock"
	appStrategy "duffett/app/strategy"
	appUser "duffett/app/user"
	_ "duffett/docs"
	"duffett/middleware"
	"duffett/pkg"
)

func init() {
	pkg.InitConfig("conf.yaml")
	pkg.InitLog(pkg.DebugLevel)
	pkg.InitDB()
	pkg.InitJwt()
}

func main() {
	router := gin.Default()

	// 注册 swagger
	// url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 配置路由
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/login", appUser.Login)
			user.POST("/register", appUser.Register)
			user.GET("/testJwt", middleware.JWTAuth(), appUser.TestJwt)
		}

		data := v1.Group("/data").Use(middleware.JWTAuth())
		{
			data.POST("/tushare", appData.Tushare)
		}

		monitor := v1.Group("/monitor")
		{
			monitor.GET("/ws", appMonitor.WS)
		}

		stock := v1.Group("/stock").Use(middleware.JWTAuth())
		{
			stock.GET("/getMonitoringStocks", appStock.GetMonitoringStocks)
			stock.GET("/getStocks", appStock.GetStocks)
		}

		order := v1.Group("/order").Use(middleware.JWTAuth())
		{
			order.POST("/getOrders", appOrder.GetOrders)
		}

		strategy := v1.Group("/strategy").Use(middleware.JWTAuth())
		{
			strategy.GET("/getStrategies", appStrategy.GetStrategies)
			strategy.POST("/create", appStrategy.Create)
			strategy.POST("/update", appStrategy.Update)
			strategy.POST("/delete", appStrategy.Delete)
			strategy.POST("/testExec", appStrategy.TestExec)
		}
	}

	// 运行
	router.Run(pkg.Conf.Addr)
}
