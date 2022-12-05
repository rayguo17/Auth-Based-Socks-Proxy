package api

import (
	"encoding/json"
	"github.com/rayguo17/go-socks/user"
	"github.com/rayguo17/go-socks/util"
	"log"
	"time"
)

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
