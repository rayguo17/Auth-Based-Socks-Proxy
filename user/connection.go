/*
	Package Connections
	We store each tcp listener accepted connection as connections
*/
package user

import (
	"log"
	"net"
)

const (
	AuthDone int = 1
	CmdRecv      = 2
	Working      = 3
	Dead         = 4
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
}

//could be manually killed or by closing the socket.
func (acpCon *AcpCon) ManualClose() {
	acpCon.manualCloseChan <- true
}

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

}
func NewCon(id string, conn net.Conn, username string, passwd string, status int) AcpCon {
	return AcpCon{
		id:         id,
		AuthChan:   make(chan bool),
		readChan:   make(chan []byte),
		writeChan:  make(chan []byte),
		conn:       conn,
		bytesCount: 0,
		owner:      nil,
		username:   username,
		passwd:     passwd,
		status:     status,
	}
}
