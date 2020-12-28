package monitor

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/99MyCql/duffett/pkg"
)

var (
	userTsMonitor = make(map[string]map[string]*monitor) // 每个用户都可监视多个股票
)

type monitorReqData struct {
	Op           string `json:"op"`            // 操作类型
	TsCode       string `json:"ts_code"`       // 股票代码
	StrategyName string `json:"strategy_name"` // 策略名字
	MonitorFreq  int64  `json:"monitor_freq"`  // 监听频率，以秒为单位
}

// @Summary StartMonitor
// @Tags Monitor
// @Accept json
// @Router /api/v1/monitor/ws [ws]
func WS(c *gin.Context) {
	// 切换为 websocket 连接
	ws, err := pkg.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusOK, pkg.ServerErr("升级websocket失败"))
		return
	}
	defer ws.Close()

	// 先获取 token
	_, tokenByte, err := ws.ReadMessage()
	if err != nil {
		log.Print(err)
		ws.WriteJSON(pkg.ServerErr("read message from websocket fail"))
		return
	}
	// 解析 token
	myClaims, err := pkg.ParseToken(string(tokenByte))
	if err != nil {
		log.Print(err)
		ws.WriteJSON(pkg.ServerErr("解析 token 失败"))
		return
	}
	// 判断 token 是否过期
	if myClaims.ExpiresAt <= time.Now().Unix() {
		log.Print(err)
		ws.WriteJSON(pkg.ServerErr("token 已过期，请重新登录"))
		return
	}
	ws.WriteJSON(pkg.Suc("websocket token 验证成功"))

	// 为每个用户创建 map[string]*monitor （如果已存在则跳过）
	username := myClaims.Username
	if _, ok := userTsMonitor[username]; !ok {
		userTsMonitor[username] = make(map[string]*monitor)
	}

	// 为当前用户下的每个运行中监听器设置ws值
	setWS(username, ws)

	for {
		// 由于WebSocket一旦连接，便可以保持长时间通讯，则该接口函数可以一直运行下去，直到连接断开
		_, dataByte, err := ws.ReadMessage()
		if err != nil {
			log.Print(err)
			ws.WriteJSON(pkg.ServerErr("read message from websocket fail"))
			break
		}
		var data monitorReqData
		if err := json.Unmarshal(dataByte, &data); err != nil {
			log.Print(err)
			ws.WriteJSON(pkg.ServerErr("json解析失败"))
			continue
		}

		if data.Op == "startMonitor" {
			// 创建并启动一个监听器
			if _, ok := userTsMonitor[username][data.TsCode]; ok {
				ws.WriteJSON(pkg.ClientErr("该监听器正在运行"))
				continue
			}
			userTsMonitor[username][data.TsCode] = newMonitor(
				username, data.TsCode, data.StrategyName, data.MonitorFreq, ws)
			if userTsMonitor[username][data.TsCode] == nil {
				continue
			}
			go userTsMonitor[username][data.TsCode].monitoring()
			ws.WriteJSON(pkg.Suc("启动监听成功"))
		} else if data.Op == "finishMonitor" {
			userTsMonitor[username][data.TsCode].finish()
			delete(userTsMonitor[username], data.TsCode)
			ws.WriteJSON(pkg.Suc("结束监听成功"))
		} else {
			ws.WriteJSON(pkg.ClientErr("op 字段未匹配"))
		}
	}
}

func setWS(username string, ws *websocket.Conn) {
	for _, v := range userTsMonitor[username] {
		v.ws = ws
	}
}
