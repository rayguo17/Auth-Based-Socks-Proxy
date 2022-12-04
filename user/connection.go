/*
	Package Connections
	We store each tcp listener accepted connection as connections
*/
package user

import (
	"errors"
	"fmt"
	"github.com/rayguo17/go-socks/util"
	"log"
	"net"
	"time"
)

const (
	AuthDone int = 1
	CmdRecv      = 2
	Working      = 3
	Dead         = 4
)

//COMMAND TYPE
const (
	Connect int = 1
	Bind        = 2
	UdpAsso     = 3
)
const BUFMAX = 4096

type AcpCon struct {
	id              string // identifier ("address:port")
	owner           *User
	bytesCount      int
	AuthChan        chan bool
	conn            net.Conn
	username        string //should be abstract to auth
	passwd          string
	endChan         chan bool
	manualCloseChan chan bool // trigger from outer to manually close
	cmdType         int
	cmdClosedChan   chan bool //tell cmd to close
	acpDelChan      chan bool //um tell acpCon this is deleted.
	cmdExecutor     Cmd
	status          int
}

func NewCon(id string, conn net.Conn, username string, passwd string) AcpCon {
	return AcpCon{
		id:         id,
		AuthChan:   make(chan bool),
		acpDelChan: make(chan bool),
		endChan:    make(chan bool),
		conn:       conn,
		bytesCount: 0,
		owner:      nil,
		//should have a auth interface...
		username: username,
		passwd:   passwd,
		status:   AuthDone,
	}
}

func (acpCon *AcpCon) EndCommand() {
	if acpCon.status == Dead {
		fmt.Println("already dead")
		return
	}
	fmt.Println("Still alive ending")
	acpCon.endChan <- true
}

//could be manually killed or by closing the socket.
func (acpCon *AcpCon) ManualClose() {
	//fmt.Println("manual close executed")
	if acpCon.status == Dead {
		return
	}
	acpCon.manualCloseChan <- true
}

//detail routine should be maintain by sub type
func (acpCon *AcpCon) MainRoutine() {
	fmt.Println("Acp MainRoutine running")
	count := 0
	for {
		count++
		//fmt.Printf("count: %d\n", count)
		select {
		case <-acpCon.manualCloseChan:
			acpCon.handleClose()
			return
		case <-acpCon.endChan:
			fmt.Println("ending acpCon")
			acpCon.handleEnd()
			return
		}
	}
	//fmt.Println("Main routine dead")
}
func (acpCon *AcpCon) handleEnd() {
	acpCon.status = Dead
	//fmt.Println("handling end in acpCon")
	go UM.DelCon(acpCon.username + "|" + acpCon.id)
	//should add a timeout don't wait forever
	select {
	case <-acpCon.acpDelChan:
		log.Println("acp con delete success")
	case <-time.After(5 * time.Second):
		log.Println("delete time out quiting anyway")
	}
	return
}

func (acpCon *AcpCon) handleClose() {
	acpCon.status = Dead
	//TODO:
	//end sub routine (if exist), delete from user manager.
	//fmt.Println("handling close")
	acpCon.conn.Close()
	if acpCon.cmdExecutor != nil {
		go acpCon.cmdExecutor.Close()
		//fmt.Println("acpCon Executing")
		//
		select {
		case <-acpCon.cmdClosedChan:
			log.Println("sub executor delete success")
		case <-time.After(5 * time.Second):
			log.Println("subExec delete timeout, quiting anyway...")
		}
	}
	//fmt.Println("deleting from um")
	//delete from  um
	go UM.DelCon(acpCon.username + "|" + acpCon.id)
	//should add a timeout don't wait forever
	select {
	case <-acpCon.acpDelChan:
		log.Println("acp con delete success")
	case <-time.After(5 * time.Second):
		log.Println("delete time out quiting anyway")
	}
	return
}
func (acpCon *AcpCon) ConnectCmd(addr util.Address) error {
	//try to dial create sub routine then return.
	//fmt.Println(addr)
	//TODO: should check the ruleset of the addr first, maybe it is not allowed
	addrStr := addr.String()
	conn, err := net.DialTimeout("tcp", addrStr, time.Second*5)
	if err != nil {
		return err
	}
	//success, create executor
	acpCon.cmdClosedChan = make(chan bool)
	connectExe := NewConExe(acpCon.cmdClosedChan, conn, addr, acpCon)
	acpCon.cmdExecutor = connectExe
	fmt.Println("command execute")
	return nil
}
func (acpCon *AcpCon) ExecuteBegin() error {
	//should check everything before begin
	go acpCon.MainRoutine()

	return acpCon.cmdExecutor.Start()
}
func (acpCon *AcpCon) CmdResponse() ([]byte, error) {
	if acpCon.cmdExecutor == nil {
		return nil, errors.New("command Executor has not been initialize")
	}
	return acpCon.cmdExecutor.FormByte(), nil
}
