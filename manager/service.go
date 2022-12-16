package manager

import (
	"encoding/json"
	"github.com/rayguo17/go-socks/manager/common"
	"github.com/rayguo17/go-socks/util"
)

//act as manager service... return light structure.

type Display struct {
	Username         string
	UploadTraffic    string
	DownloadTraffic  string
	Enable           string
	LastSeen         string
	Route            string
	ActiveConnection string
	TotalConnection  string
	Black            string
}
type GetAllUserWrap struct {
	informChan chan []*Display
}

func NewGAUWrap(informChan chan []*Display) *GetAllUserWrap {
	return &GetAllUserWrap{
		informChan: informChan,
	}
}

func (um *Manager) GetAllUser(wrap *GetAllUserWrap) {
	go func() {
		um.GetUserChannel <- wrap
	}()

}
func (um *Manager) HandleGetAllUser(wrap *GetAllUserWrap) {
	res := make([]*Display, 0, len(um.Users))
	for _, user := range um.Users {

		userConfig := &Display{
			Username:         user.GetName(),
			UploadTraffic:    user.GetUpTraffic(),
			DownloadTraffic:  user.GetDownTraffic(),
			Enable:           user.GetEnable(),
			Route:            user.GetRoute(),
			LastSeen:         user.GetLastSeen(),
			ActiveConnection: user.GetActCon(),
			TotalConnection:  user.GetTotalCon(),
			Black:            user.GetBlack(),
		}
		res = append(res, userConfig)
	}
	wrap.informChan <- res
}
func (um *Manager) AddUser(wrap *common.UserWrap) {
	go func() {
		um.AddUserChannel <- wrap
	}()
}
func (um *Manager) handleAddUser(wrap *common.UserWrap) {
	_, _, err := um.findUserByName(wrap.User.GetName())
	if err == nil {
		res := util.NewResponse(-1, "manager already exist", nil)
		wrap.InformChan <- res
		return
	}
	um.Users = append(um.Users, wrap.User)
	userConfig := &Display{
		Username: wrap.User.GetName(),
		LastSeen: wrap.User.GetLastSeen(),
		Black:    wrap.User.GetBlack(),
	}
	data, err := json.Marshal(userConfig)
	if err != nil {
		res := util.NewResponse(-1, err.Error(), nil)
		wrap.InformChan <- res
		return
	}
	res := util.NewResponse(0, "", data)
	wrap.InformChan <- res
	return
}
func (um *Manager) DelUser(wrap *common.NameWrap) {
	go func() {
		um.DelUserChannel <- wrap
	}()
}
func (um *Manager) handleDelUser(wrap *common.NameWrap) {
	//delete all acp connections. then del manager.
	//should set it to diable, and then wait until all its connection end. delete it.
	i, user, err := um.findUserByName(wrap.Username)
	if err != nil {
		wrap.InformChan <- util.NewResponse(-1, err.Error(), nil)
		return
	}
	user.SetDeleted()
	um.CheckDeletedUser(user, i)
	wrap.InformChan <- util.NewResponse(0, "", nil)
	return
}
