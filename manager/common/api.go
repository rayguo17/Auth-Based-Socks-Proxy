package common

import (
	"github.com/rayguo17/go-socks/manager/user"
	"github.com/rayguo17/go-socks/util"
)

type DCWrap struct {
	Id         string
	InformChan chan *util.Response
}

type UserWrap struct {
	User       *user.User
	InformChan chan *util.Response
}

func NUWrap(u *user.User, c chan *util.Response) *UserWrap {
	return &UserWrap{
		User:       u,
		InformChan: c,
	}
}

type NameWrap struct {
	Username   string
	InformChan chan *util.Response
}

func NNWrap(name string, c chan *util.Response) *NameWrap {
	return &NameWrap{
		name,
		c,
	}
}

type ChangePwdWrap struct {
	Username   string
	NewPasswd  string
	InformChan chan *util.Response
}

type TrafficReqWrap struct {
	Reset   bool
	ResChan chan *util.Response
}

type RulesetModWrap struct {
	Username   string
	rule       string
	Add        bool //true:add | false:delete
	InformChan chan *util.Response
}

type UploadTrafficWrap struct {
	Username   string
	Up         bool //true:up | false:down
	Count      int64
	InformChan chan *util.Response
}

type CheckRulesetWrap struct {
	Username   string
	DstAddr    string //when storing rule set, should have both domain name and ip address otherwise, we only check for the exact same address. (need to be both ipv4/domain/ipv6 yet, all the character should be the same.
	InformChan chan *util.Response
}
