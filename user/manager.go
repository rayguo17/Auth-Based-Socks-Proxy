package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rayguo17/go-socks/util"
	"github.com/rayguo17/go-socks/util/logger"
	"os"
	"strings"
	"time"
)

type Manager struct {
	Users                 []*User
	ActiveConnectionCount int
	TotalConnectionCount  int
	GetUserChannel        chan *GetAllUserWrap
	AddUserChannel        chan *UserWrap
	DelUserChannel        chan *NameWrap
	CmdChannel            chan string                   //for read write
	AcpConnections        map[string]map[string]*AcpCon //hash map, each user have a acp connections list.
	PrintUserChannel      chan bool
	AddConChannel         chan *AcpCon //after assertion need to notify, once notify done, can be continued.
	DelConChannel         chan string  // use string to delete
	ChangePwdChannel      chan *ChangePwdWrap
	TrafficReqChannel     chan *TrafficReqWrap
	RulesetModChannel     chan *RulesetModWrap
	UploadTrafficChannel  chan *UploadTrafficWrap
	CheckRulesetChannel   chan *CheckRulesetWrap
	PrintConnChannel      chan bool
}

var UM Manager
var filePath string = "./user.json"

func (um *Manager) ListUsers() {
	um.PrintUserChannel <- true
}
func (um *Manager) MainRoutine(startChan chan bool) {
	startChan <- true
	for {
		select {
		case command := <-um.CmdChannel:
			um.handleCommand(command)
		case acpCon := <-um.AddConChannel:
			um.handleAddCon(acpCon)
		case id := <-um.DelConChannel:
			um.handleDelCon(id)
		case wrap := <-um.GetUserChannel:
			um.HandleGetAllUser(wrap)
		case userWrap := <-um.AddUserChannel:
			um.handleAddUser(userWrap)
		case nameWrap := <-um.DelUserChannel:
			um.handleDelUser(nameWrap)
		case wrap := <-um.ChangePwdChannel:
			um.handleChangePwd(wrap)
		case wrap := <-um.TrafficReqChannel:
			um.handleTrafficReq(wrap)
		case wrap := <-um.RulesetModChannel:
			um.handleRulesetMod(wrap)
		case wrap := <-um.UploadTrafficChannel:
			um.handleTrafficUpload(wrap)
		case wrap := <-um.CheckRulesetChannel:
			um.handleCheckRuleset(wrap)
		case <-um.PrintConnChannel:
			um.handlePrintConn()
		case <-um.PrintUserChannel:
			um.handleUserPrintConn()
		}
	}

}
func (um *Manager) handleUserPrintConn() {
	up := util.NewUserPrinter(len(um.Users))
	for _, user := range um.Users {
		entry := util.NewUserEntry(user.GetName(), user.GetUpTraffic(), user.GetDownTraffic(), util.BoolToString(user.IsEnabled()), util.BoolToString(user.IsDeleted()), user.GetLastSeen(), user.GetActCon(), user.GetTotalCon())
		up.AddUser(entry)
	}
	up.PrintStatus()
}
func (um *Manager) handlePrintConn() {
	cp := util.NewConnPrinter(um.ActiveConnectionCount, um.TotalConnectionCount)
	for username, acpMap := range um.AcpConnections {
		for id, acp := range acpMap {

			conn := util.NewConnection(id, username, acp.RemoteAddress(), EXECSTATUS[acp.ExecutorStatus()], ACPSTATUSMAP[acp.status], CMDMap[acp.cmdType])
			cp.AddCon(conn)
		}
	}
	cp.PrintStatus()
}

func (um *Manager) handleCheckRuleset(wrap *CheckRulesetWrap) {

	_, user, err := um.findUserByName(wrap.Username)
	if err != nil {
		//log.Println("find user fail")
		res := util.NewResponse(-1, err.Error(), nil)
		wrap.informChan <- res
		return
	}
	if user.Deleted || !user.Enable {
		res := util.NewResponse(-1, "user deleted or not enabled", nil)
		wrap.informChan <- res
		return
	}
	if user.Access.Black {
		//log.Println("finding in black list")
		//check black list
		for _, v := range user.Access.BlackList {
			if v == wrap.DstAddr {
				res := util.NewResponse(-1, "dst addr in blacklist, access denied", nil)
				wrap.informChan <- res
				return
			}
		}

		res := util.NewResponse(0, "", nil)
		wrap.informChan <- res
		return
	} else {
		//log.Println("finding in white list")
		for _, v := range user.Access.WhiteList {
			if v == wrap.DstAddr {
				res := util.NewResponse(0, "", nil)
				wrap.informChan <- res
				return
			}
		}
		res := util.NewResponse(-1, "dst addr not in whitelist, access denied", nil)
		wrap.informChan <- res
		return
	}
	//log.Println("ending check rule set")
}
func (um *Manager) handleTrafficUpload(wrap *UploadTrafficWrap) {
	_, user, err := um.findUserByName(wrap.Username)
	if err != nil {
		res := util.NewResponse(-1, err.Error(), nil)
		wrap.informChan <- res
		return
	}
	if wrap.up {
		user.UplinkTraffic += wrap.count
	} else {
		user.DownLinkTraffic += wrap.count
	}
	user.lastSeen = time.Now()
	res := util.NewResponse(0, "", nil)
	wrap.informChan <- res

	return
}
func (um *Manager) handleRulesetMod(wrap *RulesetModWrap) {

}
func (um *Manager) handleTrafficReq(wrap *TrafficReqWrap) {

}
func (um *Manager) handleChangePwd(wrap *ChangePwdWrap) {

}

func (um *Manager) CheckRuleset(wrap *CheckRulesetWrap) {
	um.CheckRulesetChannel <- wrap
}
func (um *Manager) UploadTraffic(wrap *UploadTrafficWrap) {
	um.UploadTrafficChannel <- wrap
}
func (um *Manager) DelCon(id string) {
	//fmt.Println("Del con received")
	um.DelConChannel <- id
}
func (um *Manager) handleDelCon(id string) {
	//username|ip:port
	//fmt.Println("delete handling")
	//fmt.Println(id)
	idArr := strings.Split(id, "|")
	if user, ok := um.AcpConnections[idArr[0]]; ok {
		if cons, ok := user[idArr[1]]; ok {
			delete(user, idArr[1])
			index, user, err := um.findUserByName(idArr[0])
			if err != nil {
				logger.Debug.Fatal("user not found when deleting")
			}
			user.SubActiveConn()
			um.ActiveConnectionCount -= 1
			cons.acpDelChan <- true
			um.CheckDeletedUser(user, index)
		} else {
			//connection not found
			logger.Debug.Fatal("connection not found")
			return
		}
	} else {
		logger.Debug.Fatal("user not found!")
	}

}
func (um *Manager) CheckDeletedUser(u *User, i int) {
	//only the last connection could be able to delete
	if !u.Occupied() && u.IsDeleted() {
		//delete user acpCon entry.
		delete(um.AcpConnections, u.GetName())
		//delete user TODO://persistent???
		um.removeNthUser(i) //depend on user find correctness
	}
}
func (um *Manager) removeNthUser(i int) {
	size := len(um.Users)
	um.Users[i] = um.Users[size-1]
	um.Users = um.Users[:size-1]
}
func (um *Manager) ListConn() {
	um.PrintConnChannel <- true
}
func (um *Manager) AddCon(con *AcpCon) {
	um.AddConChannel <- con
}
func (um *Manager) handleAddCon(acpCon *AcpCon) {
	_, user, err := um.findUserByName(acpCon.username)
	if err != nil {
		logger.Access.Println(acpCon.Log() + " rejected user does not exist")
		acpCon.AuthChan <- false
		return
	}
	if user.Password != acpCon.passwd {
		logger.Access.Println(acpCon.Log() + " rejected password incorrect")
		acpCon.AuthChan <- false
		return
	}
	if user.Deleted || !user.Enable {
		logger.Access.Println(acpCon.Log() + " rejected user deleted or not enabled")
		acpCon.AuthChan <- false
		return
	}
	acpCon.owner = user
	if _, ok := um.AcpConnections[user.GetName()]; !ok {
		um.AcpConnections[user.GetName()] = make(map[string]*AcpCon, 0)
	}
	um.AcpConnections[user.GetName()][acpCon.id] = acpCon
	um.ActiveConnectionCount += 1
	um.TotalConnectionCount += 1
	user.AddConCount()
	acpCon.AuthChan <- true

}
func (um *Manager) findUserByName(uname string) (int, *User, error) {
	for i := 0; i < len(um.Users); i++ {
		if um.Users[i].Username == uname {
			return i, um.Users[i], nil
		}
	}
	return -1, nil, errors.New("user not found")
}
func (um *Manager) handleCommand(cmd string) {
	fmt.Println("cmd received:", cmd)
}
func (um *Manager) Initialize(path string) error {
	//read from json file, then form user group
	var Users []*User

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileBytes, &Users)
	if err != nil {
		return err
	}
	//pp.Println(Users)
	UM.Users = Users
	//initialize UM channel

	UM.GetUserChannel = make(chan *GetAllUserWrap)
	UM.AddUserChannel = make(chan *UserWrap)
	UM.DelUserChannel = make(chan *NameWrap)
	UM.ChangePwdChannel = make(chan *ChangePwdWrap)
	UM.TrafficReqChannel = make(chan *TrafficReqWrap)
	UM.RulesetModChannel = make(chan *RulesetModWrap)
	UM.UploadTrafficChannel = make(chan *UploadTrafficWrap)
	UM.CheckRulesetChannel = make(chan *CheckRulesetWrap)
	UM.PrintConnChannel = make(chan bool)
	UM.PrintUserChannel = make(chan bool)

	UM.AddConChannel = make(chan *AcpCon)
	UM.DelConChannel = make(chan string)
	UM.CmdChannel = make(chan string)
	UM.AcpConnections = make(map[string]map[string]*AcpCon)
	//pp.Println(UM)
	return nil
}
