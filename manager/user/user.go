package user

import (
	"github.com/rayguo17/go-socks/util"
	"github.com/rayguo17/go-socks/util/logger"
	"strconv"
	"strings"
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
	Route           Route `json:"Route"`
}
type Route struct {
	Type   string `json:"type"`   // Direct | Remote
	Remote string `json:"remote"` // ip:port
}
type Access struct {
	Black     bool     `json:"black"`
	BlackList []string `json:"black_list"` //support ipv4/ipv6/domain need to identify different. when matching with DstAddr, should handle carefully.
	WhiteList []string `json:"white_list"`
}

func (u *User) IsRemote() bool {
	if u.Route.Type == "" {
		return false
	}
	if u.Route.Type == "Direct" {
		return false
	}
	if u.Route.Type == "Remote" {
		return true
	}
	logger.Debug.Println("is remote status unrecognized, using default")
	return false
}
func (u *User) GetRemote() (util.Address, error) {
	//should support domanin name
	remoteArr := strings.Split(u.Route.Remote, ":")

	return util.Ipv4FromString(remoteArr[0], remoteArr[1])

}

func (u *User) GetActCon() string {
	return strconv.Itoa(u.ActiveConn)
}
func (u *User) GetTotalCon() string {
	return strconv.Itoa(u.TotalConn)
}
func (u *User) SetLastSeen(time2 time.Time) {
	u.lastSeen = time2
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
func (u *User) GetUpTraffic() string {
	return strconv.FormatInt(u.UplinkTraffic, 10)
}
func (u *User) GetDownTraffic() string {
	return strconv.FormatInt(u.DownLinkTraffic, 10)
}
func (u *User) IsEnabled() bool {
	return u.Enable
}
func (u *User) IsDeleted() bool {
	return u.Deleted
}
func (u *User) Occupied() bool {
	if u.ActiveConn != 0 {
		return true
	}
	return false
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
