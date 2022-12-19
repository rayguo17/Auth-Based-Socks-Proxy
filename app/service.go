package app

import (
	"encoding/json"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/manager"
	"github.com/rayguo17/go-socks/manager/common"
	"github.com/rayguo17/go-socks/manager/user"
	"github.com/rayguo17/go-socks/util"
	"github.com/rayguo17/go-socks/util/logger"
	"time"
)

func GetAllUser() *util.AppResponse {
	informChan := make(chan []*manager.Display)
	getAllUserWrap := manager.NewGAUWrap(informChan)
	manager.UM.GetAllUser(getAllUserWrap)
	select {
	case userDisplay := <-informChan:
		res := util.NewAppResponse(0, "", userDisplay)
		return res
	case <-time.After(time.Second * 5):
		logger.Debug.Println("get all user timeout")
		res := util.NewAppResponse(-1, "get user timeout", nil)
		return res
	}
}
func AddUserHandler(newUser interface{}) *util.Response {
	c := make(chan *util.Response)
	data, err := json.Marshal(newUser)
	if err != nil {
		return util.NewResponse(-1, err.Error(), nil)
	}
	var u user.User
	err = json.Unmarshal(data, &u)
	if err != nil {
		return util.NewResponse(-1, err.Error(), nil)
	}
	uw := common.NUWrap(&u, c)
	manager.UM.AddUser(uw)
	select {
	case res := <-c:
		return res
	case <-time.After(time.Second * 5):
		return util.NewResponse(-1, "get manager timeout", nil)
	}
}
func DelUserHandler(name string) *util.Response {
	c := make(chan *util.Response)
	nwrap := common.NNWrap(name, c)
	manager.UM.DelUser(nwrap)
	select {
	case res := <-c:
		return res
	case <-time.After(time.Second * 5):
		return util.NewResponse(-1, "get manager timeout", nil)
	}
}
func ConfigHandler() *Response {
	return NewAppResponse(0, "", config.SystemConfig)
}
