package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rayguo17/go-socks/manager/common"
	"github.com/rayguo17/go-socks/manager/connection"
	"github.com/rayguo17/go-socks/manager/user"
	"github.com/rayguo17/go-socks/util"
	"github.com/rayguo17/go-socks/util/logger"
	"os"
	"strings"
	"time"
)

type Manager struct {
	Users                 []*user.User
	ActiveConnectionCount int
	TotalConnectionCount  int
	GetUserChannel        chan *GetAllUserWrap
	AddUserChannel        chan *common.UserWrap
	DelUserChannel        chan *common.NameWrap
	CmdChannel            chan string                              //for read write
	AcpConnections        map[string]map[string]*connection.AcpCon //hash map, each manager have a acp connections list.
	PrintUserChannel      chan bool
	AddConChannel         chan *connection.AcpCon //after assertion need to notify, once notify done, can be continued.
	DelConChannel         chan *common.DCWrap     // use string to delete
	ChangePwdChannel      chan *common.ChangePwdWrap
	TrafficReqChannel     chan *common.TrafficReqWrap
	RulesetModChannel     chan *common.RulesetModWrap
	UploadTrafficChannel  chan *common.UploadTrafficWrap
	CheckRulesetChannel   chan *common.CheckRulesetWrap
	PrintConnChannel      chan bool
}

var UM Manager
var filePath string = "./manager.json"

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
func (um *Manager) GetConCommunicator() *common.Communicator {
	return common.NewCommunicator(um.CheckRulesetChannel, um.UploadTrafficChannel, um.DelConChannel)
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

			conn := util.NewConnection(id, username, acp.RemoteAddress(), connection.EXECSTATUS[acp.ExecutorStatus()], connection.ACPSTATUSMAP[acp.GetStatus()], connection.CMDMap[acp.GetCmdType()])
			cp.AddCon(conn)
		}
	}
	cp.PrintStatus()
}

func (um *Manager) handleCheckRuleset(wrap *common.CheckRulesetWrap) {

	_, user, err := um.findUserByName(wrap.Username)
	if err != nil {
		//log.Println("find manager fail")
		res := util.NewResponse(-1, err.Error(), nil)
		wrap.InformChan <- res
		return
	}
	if user.Deleted || !user.Enable {
		res := util.NewResponse(-1, "manager deleted or not enabled", nil)
		wrap.InformChan <- res
		return
	}
	if user.Access.Black {
		//log.Println("finding in black list")
		//check black list
		for _, v := range user.Access.BlackList {
			if v == wrap.DstAddr {
				res := util.NewResponse(-1, "dst addr in blacklist, access denied", nil)
				wrap.InformChan <- res
				return
			}
		}

		res := util.NewResponse(0, "", nil)
		wrap.InformChan <- res
		return
	} else {
		//log.Println("finding in white list")
		for _, v := range user.Access.WhiteList {
			if v == wrap.DstAddr {
				res := util.NewResponse(0, "", nil)
				wrap.InformChan <- res
				return
			}
		}
		res := util.NewResponse(-1, "dst addr not in whitelist, access denied", nil)
		wrap.InformChan <- res
		return
	}
	//log.Println("ending check rule set")
}
func (um *Manager) handleTrafficUpload(wrap *common.UploadTrafficWrap) {
	_, user, err := um.findUserByName(wrap.Username)
	if err != nil {
		res := util.NewResponse(-1, err.Error(), nil)
		wrap.InformChan <- res
		return
	}
	if wrap.Up {
		user.UplinkTraffic += wrap.Count
	} else {
		user.DownLinkTraffic += wrap.Count
	}
	user.SetLastSeen(time.Now())
	res := util.NewResponse(0, "", nil)
	wrap.InformChan <- res

	return
}
func (um *Manager) handleRulesetMod(wrap *common.RulesetModWrap) {

}
func (um *Manager) handleTrafficReq(wrap *common.TrafficReqWrap) {

}
func (um *Manager) handleChangePwd(wrap *common.ChangePwdWrap) {

}

func (um *Manager) CheckRuleset(wrap *common.CheckRulesetWrap) {
	um.CheckRulesetChannel <- wrap
}
func (um *Manager) UploadTraffic(wrap *common.UploadTrafficWrap) {
	go func() {
		um.UploadTrafficChannel <- wrap
	}()
}
func (um *Manager) DelCon(id *common.DCWrap) {
	//fmt.Println("Del con received")
	um.DelConChannel <- id
}
func (um *Manager) handleDelCon(wrap *common.DCWrap) {
	//username|ip:port
	//fmt.Println("delete handling")
	//fmt.Println(id)
	id := wrap.Id
	idArr := strings.Split(id, "|")
	if user, ok := um.AcpConnections[idArr[0]]; ok {
		if _, ok := user[idArr[1]]; ok {
			delete(user, idArr[1])
			index, user, err := um.findUserByName(idArr[0])
			if err != nil {
				logger.Debug.Fatal("manager not found when deleting")
			}
			user.SubActiveConn()
			um.ActiveConnectionCount -= 1
			res := util.NewResponse(0, "", nil)
			wrap.InformChan <- res
			um.CheckDeletedUser(user, index)
		} else {
			//connection not found
			logger.Debug.Fatal("connection not found")
			return
		}
	} else {
		logger.Debug.Fatal("manager not found!")
	}

}
func (um *Manager) CheckDeletedUser(u *user.User, i int) {
	//only the last connection could be able to delete
	if !u.Occupied() && u.IsDeleted() {
		//delete manager acpCon entry.
		delete(um.AcpConnections, u.GetName())
		//delete manager TODO://persistent???
		um.removeNthUser(i) //depend on manager find correctness
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
func (um *Manager) AddCon(con *connection.AcpCon) {
	um.AddConChannel <- con
}
func (um *Manager) handleAddCon(acpCon *connection.AcpCon) {
	_, user, err := um.findUserByName(acpCon.GetName())
	if err != nil {
		logger.Access.Println(acpCon.Log() + " rejected manager does not exist")
		acpCon.AuthChan <- false
		return
	}
	if user.Password != acpCon.GetPasswd() {
		logger.Access.Println(acpCon.Log() + " rejected password incorrect")
		acpCon.AuthChan <- false
		return
	}
	if user.Deleted || !user.Enable {
		logger.Access.Println(acpCon.Log() + " rejected manager deleted or not enabled")
		acpCon.AuthChan <- false
		return
	}
	if user.IsRemote() {
		addr, err := user.GetRemote()
		if err != nil {
			acpCon.SetRemote(false, nil)
		} else {
			acpCon.SetRemote(true, addr)
		}
	}
	if _, ok := um.AcpConnections[user.GetName()]; !ok {
		um.AcpConnections[user.GetName()] = make(map[string]*connection.AcpCon, 0)
	}
	um.AcpConnections[user.GetName()][acpCon.GetID()] = acpCon
	um.ActiveConnectionCount += 1
	um.TotalConnectionCount += 1
	user.AddConCount()
	acpCon.AuthChan <- true

}
func (um *Manager) findUserByName(uname string) (int, *user.User, error) {
	for i := 0; i < len(um.Users); i++ {
		if um.Users[i].Username == uname {
			return i, um.Users[i], nil
		}
	}
	return -1, nil, errors.New("manager not found")
}
func (um *Manager) handleCommand(cmd string) {
	fmt.Println("cmd received:", cmd)
}
func (um *Manager) Initialize(path string) error {
	//read from json file, then form manager group
	var Users []*user.User

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
	UM.AddUserChannel = make(chan *common.UserWrap)
	UM.DelUserChannel = make(chan *common.NameWrap)
	UM.ChangePwdChannel = make(chan *common.ChangePwdWrap)
	UM.TrafficReqChannel = make(chan *common.TrafficReqWrap)
	UM.RulesetModChannel = make(chan *common.RulesetModWrap)
	UM.UploadTrafficChannel = make(chan *common.UploadTrafficWrap)
	UM.CheckRulesetChannel = make(chan *common.CheckRulesetWrap)
	UM.PrintConnChannel = make(chan bool)
	UM.PrintUserChannel = make(chan bool)

	UM.AddConChannel = make(chan *connection.AcpCon)
	UM.DelConChannel = make(chan *common.DCWrap)
	UM.CmdChannel = make(chan string)
	UM.AcpConnections = make(map[string]map[string]*connection.AcpCon)
	//pp.Println(UM)
	return nil
}
