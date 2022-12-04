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
	Username string `json:"username"`
	Password string `json:"password"`
}

//encode cmd into string. parse
type Manager struct {
	Users          []*User
	CmdChannel     chan string                   //for read write
	AcpConnections map[string]map[string]*AcpCon //hash map, each user have a acp connections list.
	printChannel   chan bool
	AddConChannel  chan *AcpCon //after assertion need to notify, once notify done, can be continued.
	DelConChannel  chan string  // use string to delete
	AddUserChannel chan *User
	DelUserChannel chan string //use username to delete
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
			pp.Println(um.AcpConnections)
		case acpCon := <-um.AddConChannel:
			um.handleAdd(acpCon)
		case id := <-um.DelConChannel:
			um.handleDel(id)
		}
	}

}
func (um *Manager) DelCon(id string) {
	//fmt.Println("Del con received")
	um.DelConChannel <- id
}
func (um *Manager) handleDel(id string) {
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
func (um *Manager) handleAdd(acpCon *AcpCon) {
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
	UM.Users = Users
	//initialize UM channel
	UM.DelConChannel = make(chan string)
	UM.AddConChannel = make(chan *AcpCon)
	UM.printChannel = make(chan bool)
	UM.CmdChannel = make(chan string)
	UM.AcpConnections = make(map[string]map[string]*AcpCon)
	//pp.Println(UM)
}
