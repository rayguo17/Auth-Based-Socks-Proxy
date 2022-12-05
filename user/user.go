package user

import (
	"time"
)

//upper case so json pkg have access to field
type User struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	uplinkTraffic   int64
	downLinkTraffic int64
	lastSeen        time.Time
	Access          Access `json:"access"`
}
type Access struct {
	Black     bool     `json:"black"`
	BlackList []string `json:"black_list"` //support ipv4/ipv6/domain need to identify different. when matching with DstAddr, should handle carefully.
	WhiteList []string `json:"white_list"`
}

func (u *User) GetName() string {
	return u.Username
}
func (u *User) GetBlack() string {
	if u.Access.Black {
		return "Black"
	}
	return "White"
}

func (u *User) GetLastSeen() string {
	initTime := time.Time{}
	if u.lastSeen.Equal(initTime) {
		return ""
	} else {
		return u.lastSeen.String()
	}
}
