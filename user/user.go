package user

import (
	"encoding/json"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"github.com/rayguo17/go-socks/connections"
	"log"
	"os"
)

//upper case so json pkg have access to field
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//encode cmd into string. parse
type UserManager struct {
	Users          []*User
	CmdChannel     chan string                      //for read write
	AcpConnections map[string][]*connections.AcpCon //hash map, each user have a acp connections list.
	printChannel   chan bool
}

var Users []*User
var filePath string = "./user.json"

func (um *UserManager) ListUser() {
	um.printChannel <- true

}
func (um *UserManager) MainRoutine() {
	for {
		select {
		case command := <-um.CmdChannel:

			um.handleCommand(command)
		case <-um.printChannel:
			pp.Println(um.AcpConnections)
		}
	}

}
func (um *UserManager) handleCommand(cmd string) {
	fmt.Println("cmd received:", cmd)
}
func init() {
	//read from json file, then form user group
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
}
