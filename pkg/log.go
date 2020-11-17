package pkg

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// InitLog 初始化日志配置
func InitLog() {
	// 配置日志输出。如果未设置日志文件，则输出到控制台
	if Conf.LogPath == "" {
		log.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(Conf.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
			return
		}
		gin.DefaultWriter = file
		log.SetOutput(file)
	}
	log.SetPrefix("[duffett] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
