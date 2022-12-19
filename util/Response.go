package util

import "strconv"

type Response struct {
	errCode int
	errMsg  string
	data    []byte
}
type AppResponse struct {
	errCode int
	errMsg  string
	data    interface{}
}

func NewAppResponse(c int, msg string, data interface{}) *AppResponse {
	return &AppResponse{
		errMsg:  msg,
		errCode: c,
		data:    data,
	}
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
func (r *AppResponse) GetData() interface{} {
	return r.data
}
func (r *AppResponse) GetErrCode() int {
	return r.errCode
}
func (r *AppResponse) GetErrMsg() string {
	return r.errMsg
}
func (r *Response) String() string {
	code := strconv.Itoa(r.errCode)
	return code + ":" + r.errMsg + string(r.data)
}
