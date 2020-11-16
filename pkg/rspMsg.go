package pkg

type RspMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	SucCode  = 0
	ErrCode  = 1
	FailCode = 2
)
