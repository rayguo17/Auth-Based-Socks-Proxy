package user

import (
	"encoding/json"
	"github.com/rayguo17/go-socks/util"
)

//act as user service... return server structure.

type Display struct {
	Username string
	Black    string
	LastSeen string
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
			Username: user.GetName(),
			LastSeen: user.GetLastSeen(),
			Black:    user.GetBlack(),
		}
		res = append(res, userConfig)
	}
	wrap.informChan <- res
}
func (um *Manager) AddUser(wrap *UserWrap) {
	go func() {
		um.AddUserChannel <- wrap
	}()
}
func (um *Manager) handleAddUser(wrap *UserWrap) {
	_, _, err := um.findUserByName(wrap.user.GetName())
	if err == nil {
		res := util.NewResponse(-1, "user already exist", nil)
		wrap.informChan <- res
		return
	}
	um.Users = append(um.Users, wrap.user)
	userConfig := &Display{
		Username: wrap.user.GetName(),
		LastSeen: wrap.user.GetLastSeen(),
		Black:    wrap.user.GetBlack(),
	}
	data, err := json.Marshal(userConfig)
	if err != nil {
		res := util.NewResponse(-1, err.Error(), nil)
		wrap.informChan <- res
		return
	}
	res := util.NewResponse(0, "", data)
	wrap.informChan <- res
	return
}
func (um *Manager) DelUser(wrap *NameWrap) {
	go func() {
		um.DelUserChannel <- wrap
	}()
}
func (um *Manager) handleDelUser(wrap *NameWrap) {
	//delete all acp connections. then del user.
	//should set it to diable, and then wait until all its connection end. delete it.
	i, user, err := um.findUserByName(wrap.Username)
	if err != nil {
		wrap.informChan <- util.NewResponse(-1, err.Error(), nil)
		return
	}
	user.SetDeleted()
	um.CheckDeletedUser(user, i)
	wrap.informChan <- util.NewResponse(0, "", nil)
	return
}
