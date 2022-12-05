package util

import "strconv"

type Response struct {
	errCode int
	errMsg  string
	data    []byte
}

func NewResponse(c int, msg string, data []byte) *Response {
	return &Response{
		errCode: c,
		errMsg:  msg,
		data:    data,
	}
}
func (r *Response) GetData() []byte {
	return r.data
}
func (r *Response) GetErrCode() int {
	return r.errCode
}
func (r *Response) GetErrMsg() string {
	return r.errMsg
}
func (r *Response) String() string {
	code := strconv.Itoa(r.errCode)
	return code + ":" + r.errMsg + string(r.data)
}
