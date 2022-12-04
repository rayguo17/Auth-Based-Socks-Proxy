/*
	Package Connections
	We store each tcp listener accepted connection as connections
*/
package user

import (
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
	username        string
	passwd          string
	status          int
	readChan        chan []byte
	writeChan       chan []byte
	endChan         chan bool
	manualCloseChan chan bool
	cmdType         int
	cmdClosedChan   chan bool
	acpDelChan      chan bool
	cmdExecutor     Cmd
}

//could be manually killed or by closing the socket.
func (acpCon *AcpCon) ManualClose() {
	//fmt.Println("manual close executed")
	acpCon.handleClose()
	//acpCon.manualCloseChan <- true
}
func (acpCon *AcpCon) ExecuteCmd() {

}

//detail routine should be maintain by sub type
func (acpCon *AcpCon) MainRoutine() {
	for {
		select {
		case buf := <-acpCon.readChan:
			acpCon.handleRead(buf)
		case buf := <-acpCon.writeChan:
			acpCon.handleWrite(buf) //maybe we will need a buffer??
		case <-acpCon.manualCloseChan:
			acpCon.handleClose()
		}
	}
}
func (acpCon *AcpCon) handleWrite(buf []byte) {

}
func (acpCon *AcpCon) handleRead(buf []byte) {
	switch acpCon.status {
	case AuthDone:
		//read command

	}

}
func (acpCon *AcpCon) ReadRoutine() {
	for {
		buf := make([]byte, BUFMAX)
		readLen, err := acpCon.conn.Read(buf)
		if err != nil {
			log.Println("read routine error", err)
			acpCon.handleClose()
		}
		acpCon.readChan <- buf[:readLen]
	}
}
func (acpCon *AcpCon) handleClose() {
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
func (acpCon *AcpCon) ConnectCmd(addr string) error {
	//try to dial create sub routine then return.
	//fmt.Println(addr)
	//TODO: should check the ruleset of the addr first, maybe it is not allowed
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	//success, create executor
	acpCon.cmdClosedChan = make(chan bool)
	connectExe := NewConExe(acpCon.cmdClosedChan, conn)
	acpCon.cmdExecutor = connectExe
	return nil
}
func NewCon(id string, conn net.Conn, username string, passwd string, status int) AcpCon {
	return AcpCon{
		id:         id,
		AuthChan:   make(chan bool),
		readChan:   make(chan []byte),
		writeChan:  make(chan []byte),
		acpDelChan: make(chan bool),
		conn:       conn,
		bytesCount: 0,
		owner:      nil,
		//should have a auther interface...
		username: username,
		passwd:   passwd,
		status:   status,
	}
}
