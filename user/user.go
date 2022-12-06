package user

import (
	"time"
)

//upper case so json pkg have access to field
type User struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	UplinkTraffic   int64  `json:"uplink_traffic"`
	DownLinkTraffic int64  `json:"downLink_traffic"`
	Enable          bool   `json:"enable"`
	Deleted         bool
	lastSeen        time.Time
	Access          Access `json:"access"`
	ActiveConn      int
	TotalConn       int
}
type Access struct {
	Black     bool     `json:"black"`
	BlackList []string `json:"black_list"` //support ipv4/ipv6/domain need to identify different. when matching with DstAddr, should handle carefully.
	WhiteList []string `json:"white_list"`
}

func (u *User) DelUser() {
	u.Deleted = true
}
func (u *User) SubActiveConn() {
	u.ActiveConn -= 1
}
func (u *User) AddConCount() {
	u.ActiveConn += 1
	u.TotalConn += 1
}
func (u *User) SetDeleted() {
	u.Deleted = true
}
func (u *User) IsDeleted() bool {
	return u.Deleted
}
func (u *User) Occupied() bool {
	if u.ActiveConn != 0 {
		return false
	}
	return true
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
