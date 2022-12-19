package app

type Response struct {
	ErrCode int
	ErrMsg  string
	Data    interface{}
}

func NewAppResponse(c int, msg string, data interface{}) *Response {
	return &Response{
		ErrCode: c,
		ErrMsg:  msg,
		Data:    data,
	}
}
