package user

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
