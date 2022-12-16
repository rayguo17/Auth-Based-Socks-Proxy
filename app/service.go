package app

import (
	"github.com/rayguo17/go-socks/manager"
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
