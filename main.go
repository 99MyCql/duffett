package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"duffett/app"
	_ "duffett/docs"
	"duffett/pkg"
)

func init() {
	pkg.InitConfig()
	pkg.InitLog()
	pkg.InitDB()
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	// 访问 swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/login", app.Login)
		}
	}

	// 运行
	router.Run(pkg.Conf.Addr)
}
