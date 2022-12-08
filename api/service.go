package api

import (
	"encoding/json"
	"github.com/rayguo17/go-socks/manager"
	"github.com/rayguo17/go-socks/manager/common"
	"github.com/rayguo17/go-socks/manager/user"
	"github.com/rayguo17/go-socks/util"
	"log"
	"time"
)

type DelParams struct {
	Username string `json:"username"`
}

func GetAllUser() *util.Response {
	informChan := make(chan []*manager.Display)
	getAlluserWrap := manager.NewGAUWrap(informChan)
	manager.UM.GetAllUser(getAlluserWrap)
	select {
	case userConfigs := <-informChan:
		//pp.Println(userConfigs)
		data, err := json.Marshal(userConfigs)
		if err != nil {
			res := util.NewResponse(-1, err.Error(), nil)
			return res
		}
		return util.NewResponse(0, "", data)
	case <-time.After(time.Second * 5):
		log.Println("get all manager timeout")
		res := util.NewResponse(-1, "get manager timeout", nil)
		return res
	}

}
func AddUser(u *user.User) *util.Response {
	c := make(chan *util.Response)
	uw := common.NUWrap(u, c)
	manager.UM.AddUser(uw)
	select {
	case res := <-c:
		return res
	case <-time.After(time.Second * 5):
		return util.NewResponse(-1, "get manager timeout", nil)
	}
}
func DelUser(d *DelParams) *util.Response {
	c := make(chan *util.Response)
	nwrap := common.NNWrap(d.Username, c)
	manager.UM.DelUser(nwrap)
	select {
	case res := <-c:
		return res
	case <-time.After(time.Second * 5):
		return util.NewResponse(-1, "get manager timeout", nil)
	}
}
