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

// ClientErr 客户端错误时返回信息
func ClientErr(msg string) RspData {
	return RspData{
		Code: ClientErrCode,
		Msg:  msg,
	}
}

// ServerErr 服务端错误时返回信息
func ServerErr(msg string) RspData {
	return RspData{
		Code: ClientErrCode,
		Msg:  msg,
	}
}

// Suc 成功时返回信息
func Suc(msg string) RspData {
	return RspData{
		Code: SucCode,
		Msg:  msg,
	}
}

// SucWithData 成功时返回信息（携带数据）
func SucWithData(msg string, data interface{}) RspData {
	return RspData{
		Code: SucCode,
		Msg:  msg,
		Data: data,
	}
}
