package pkg

// RspData HTTP 返回数据体
type RspData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	SucCode       = 0 // 成功码
	ClientErrCode = 1 // 客户端错误码
	ServerErrCode = 2 // 服务器错误码
)
