package pkg

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var WsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
