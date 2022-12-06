package api

import (
	"encoding/json"
	"github.com/rayguo17/go-socks/user"
	"github.com/rayguo17/go-socks/util"
	"log"
	"time"
)

type DelParams struct {
	Username string `json:"username"`
}

func GetAllUser() *util.Response {
	informChan := make(chan []*user.Display)
	getAlluserWrap := user.NewGAUWrap(informChan)
	user.UM.GetAllUser(getAlluserWrap)
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
		log.Println("get all user timeout")
		res := util.NewResponse(-1, "get user timeout", nil)
		return res
	}

}
func AddUser(u *user.User) *util.Response {
	c := make(chan *util.Response)
	uw := user.NUWrap(u, c)
	user.UM.AddUser(uw)
	select {
	case res := <-c:
		return res
	case <-time.After(time.Second * 5):
		return util.NewResponse(-1, "get user timeout", nil)
	}
}
func DelUser(d *DelParams) *util.Response {
	c := make(chan *util.Response)
	nwrap := user.NNWrap(d.Username, c)
	user.UM.DelUser(nwrap)
	select {
	case res := <-c:
		return res
	case <-time.After(time.Second * 5):
		return util.NewResponse(-1, "get user timeout", nil)
	}
}
