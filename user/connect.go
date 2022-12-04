package user

import "net"

const (
	RUNNING int = 1
	DEAD    int = 2
)

type ConnectExecutor struct {
	targetCon net.Conn
	readChan  chan []byte
	writeChan chan []byte
	acpCon    *AcpCon
	closeChan chan bool //acpCon inform executor to close routine
	status    int
}

func NewConExe(closeChan chan bool, conn net.Conn) *ConnectExecutor {
	return &ConnectExecutor{
		closeChan: closeChan,
		readChan:  make(chan []byte),
		writeChan: make(chan []byte),
		targetCon: conn,
		status:    DEAD,
	}
}
func (ce *ConnectExecutor) Start() {

}
func (ce *ConnectExecutor) Close() {
	if ce.status != RUNNING {
		ce.targetCon.Close()
		ce.closeChan <- true
	} else {
		//close routine
	}

}
func (ce *ConnectExecutor) MainRoutine() {

}
