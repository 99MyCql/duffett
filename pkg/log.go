package pkg

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// myFormatter 实现 logrus.Formatter 接口，自定义输出格式
type myFormatter struct{}

func (f *myFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02-15:04:05")
	msg := fmt.Sprintf("[duffett] [%s] %s %s:%d %s\n", strings.ToUpper(entry.Level.String()), timestamp,
		entry.Caller.File, entry.Caller.Line, entry.Message)
	return []byte(msg), nil
}

const (
	DebugLevel = "Debug"
	InfoLevel  = "Info"
	WarnLevel  = "Warn"
	ErrorLevel = "Error"
	FatalLevel = "Fatal"
)

// InitLog 初始化日志配置
func InitLog(level string) {
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

	// 设置日志级别
	switch level {
	case DebugLevel:
		log.SetLevel(log.DebugLevel)
	case InfoLevel:
		log.SetLevel(log.InfoLevel)
	case WarnLevel:
		log.SetLevel(log.WarnLevel)
	case ErrorLevel:
		log.SetLevel(log.ErrorLevel)
	case FatalLevel:
		log.SetLevel(log.FatalLevel)
	default:
		panic("未匹配的日志级别")
	}

	// 设置在输出日志中添加文件名和方法信息
	log.SetReportCaller(true)

	// 设置自定义输出格式
	log.SetFormatter(new(myFormatter))
}
