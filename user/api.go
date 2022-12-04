package user

import "strconv"

type Response struct {
	errCode int
	errMsg  string
	data    string
}

func NewResponse(c int, msg string, data string) *Response {
	return &Response{
		errCode: c,
		errMsg:  msg,
		data:    data,
	}
}
func (r *Response) String() string {
	code := strconv.Itoa(r.errCode)
	return code + ":" + r.errMsg + r.data
}

type UserWrap struct {
	user       *User
	informChan chan *Response
}

type NameWrap struct {
	Username   string
	informChan chan *Response
}

type ChangePwdWrap struct {
	Username   string
	NewPasswd  string
	informChan chan *Response
}

type TrafficReqWrap struct {
	Reset   bool
	ResChan chan *Response
}

type RulesetModWrap struct {
	Username   string
	rule       string
	Add        bool //true:add | false:delete
	informChan chan *Response
}

type UploadTrafficWrap struct {
	Username   string
	up         bool //true:up | false:down
	count      int64
	informChan chan *Response
}

type CheckRulesetWrap struct {
	Username   string
	DstAddr    string //when storing rule set, should have both domain name and ip address otherwise, we only check for the exact same address. (need to be both ipv4/domain/ipv6 yet, all the character should be the same.
	informChan chan *Response
}
