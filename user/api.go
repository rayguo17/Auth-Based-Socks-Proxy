package user

import (
	"github.com/rayguo17/go-socks/util"
)

type UserWrap struct {
	user       *User
	informChan chan *util.Response
}

type NameWrap struct {
	Username   string
	informChan chan *util.Response
}

type ChangePwdWrap struct {
	Username   string
	NewPasswd  string
	informChan chan *util.Response
}

type TrafficReqWrap struct {
	Reset   bool
	ResChan chan *util.Response
}

type RulesetModWrap struct {
	Username   string
	rule       string
	Add        bool //true:add | false:delete
	informChan chan *util.Response
}

type UploadTrafficWrap struct {
	Username   string
	up         bool //true:up | false:down
	count      int64
	informChan chan *util.Response
}

type CheckRulesetWrap struct {
	Username   string
	DstAddr    string //when storing rule set, should have both domain name and ip address otherwise, we only check for the exact same address. (need to be both ipv4/domain/ipv6 yet, all the character should be the same.
	informChan chan *util.Response
}
