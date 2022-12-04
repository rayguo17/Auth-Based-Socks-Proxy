package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"log"
	"os"
	"strings"
)

//upper case so json pkg have access to field
type User struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	uplinkTraffic   int64
	downLinkTraffic int64
	Access          Access `json:"access"`
}
type Access struct {
	Black     bool     `json:"black"`
	BlackList []string `json:"black_list"` //support ipv4/ipv6/domain need to identify different. when matching with DstAddr, should handle carefully.
	WhiteList []string `json:"white_list"`
}

//encode cmd into string. parse
type Manager struct {
	Users                []*User
	AddUserChannel       chan *UserWrap
	DelUserChannel       chan *NameWrap
	CmdChannel           chan string                   //for read write
	AcpConnections       map[string]map[string]*AcpCon //hash map, each user have a acp connections list.
	printChannel         chan bool
	AddConChannel        chan *AcpCon //after assertion need to notify, once notify done, can be continued.
	DelConChannel        chan string  // use string to delete
	ChangePwdChannel     chan *ChangePwdWrap
	TrafficReqChannel    chan *TrafficReqWrap
	RulesetModChannel    chan *RulesetModWrap
	UploadTrafficChannel chan *UploadTrafficWrap
	CheckRulesetChannel  chan *CheckRulesetWrap
}

var UM Manager
var filePath string = "./user.json"

func (um *Manager) ListUser() {
	um.printChannel <- true
}
func (um *Manager) MainRoutine() {
	for {
		select {
		case command := <-um.CmdChannel:
			um.handleCommand(command)
		case <-um.printChannel:
			pp.Println(um.Users)
		case acpCon := <-um.AddConChannel:
			um.handleAddCon(acpCon)
		case id := <-um.DelConChannel:
			um.handleDelCon(id)
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
		}
	}

}
func (um *Manager) handleCheckRuleset(wrap *CheckRulesetWrap) {

	user, err := um.findUserByName(wrap.Username)
	if err != nil {
		//log.Println("find user fail")
		res := NewResponse(-1, err.Error(), "")
		wrap.informChan <- res
		return
	}
	if user.Access.Black {
		//log.Println("finding in black list")
		//check black list
		for _, v := range user.Access.BlackList {
			if v == wrap.DstAddr {
				res := NewResponse(-1, "dst addr in blacklist, access denied", "")
				wrap.informChan <- res
				return
			}
		}
		res := NewResponse(0, "", "")
		wrap.informChan <- res
		return
	} else {
		//log.Println("finding in white list")
		for _, v := range user.Access.WhiteList {
			if v == wrap.DstAddr {
				res := NewResponse(0, "", "")
				wrap.informChan <- res
				return
			}
		}
		res := NewResponse(-1, "dst addr not in whitelist, access denied", "")
		wrap.informChan <- res
		return
	}
	//log.Println("ending check rule set")
}
func (um *Manager) handleTrafficUpload(wrap *UploadTrafficWrap) {
	user, err := um.findUserByName(wrap.Username)
	if err != nil {
		res := NewResponse(-1, err.Error(), "")
		wrap.informChan <- res
		return
	}
	if wrap.up {
		user.uplinkTraffic += wrap.count
	} else {
		user.downLinkTraffic += wrap.count
	}
	res := NewResponse(0, "", "")
	wrap.informChan <- res
	return
}
func (um *Manager) handleRulesetMod(wrap *RulesetModWrap) {

}
func (um *Manager) handleTrafficReq(wrap *TrafficReqWrap) {

}
func (um *Manager) handleChangePwd(wrap *ChangePwdWrap) {

}
func (um *Manager) handleDelUser(wrap *NameWrap) {

}
func (um *Manager) handleAddUser(wrap *UserWrap) {

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
			cons.acpDelChan <- true
		} else {
			//connection not found
			log.Fatal("connection not found")
			return
		}
	} else {
		log.Fatal("user not found!")
	}

}
func (um *Manager) AddCon(con *AcpCon) {
	um.AddConChannel <- con
}
func (um *Manager) handleAddCon(acpCon *AcpCon) {
	user, err := um.findUserByName(acpCon.username)
	if err != nil {
		acpCon.AuthChan <- false
		return
	}
	if user.Password != acpCon.passwd {
		acpCon.AuthChan <- false
		return
	}
	acpCon.owner = user
	if _, ok := um.AcpConnections[user.Username]; !ok {
		um.AcpConnections[user.Username] = make(map[string]*AcpCon, 0)
	}
	um.AcpConnections[user.Username][acpCon.id] = acpCon
	acpCon.AuthChan <- true
}
func (um *Manager) findUserByName(uname string) (*User, error) {
	for i := 0; i < len(um.Users); i++ {
		if um.Users[i].Username == uname {
			return um.Users[i], nil
		}
	}
	return nil, errors.New("user not found")
}
func (um *Manager) handleCommand(cmd string) {
	fmt.Println("cmd received:", cmd)
}
func init() {
	//read from json file, then form user group
	var Users []*User

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("read file failed")
		return
	}
	err = json.Unmarshal(fileBytes, &Users)
	if err != nil {
		log.Fatal("unmarshal failed", err.Error())
	}
	//pp.Println(Users)
	UM.Users = Users
	//initialize UM channel

	UM.AddUserChannel = make(chan *UserWrap)
	UM.DelUserChannel = make(chan *NameWrap)
	UM.ChangePwdChannel = make(chan *ChangePwdWrap)
	UM.TrafficReqChannel = make(chan *TrafficReqWrap)
	UM.RulesetModChannel = make(chan *RulesetModWrap)
	UM.UploadTrafficChannel = make(chan *UploadTrafficWrap)
	UM.CheckRulesetChannel = make(chan *CheckRulesetWrap)

	UM.AddConChannel = make(chan *AcpCon)
	UM.DelConChannel = make(chan string)
	UM.printChannel = make(chan bool)
	UM.CmdChannel = make(chan string)
	UM.AcpConnections = make(map[string]map[string]*AcpCon)
	//pp.Println(UM)
}
