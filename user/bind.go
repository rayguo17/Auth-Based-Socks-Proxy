package user

import "net"

type BindExecutor struct {
	targetCon net.Conn
	readChan  chan []byte
	writeChan chan []byte
	acpCon    *AcpCon
	closeChan chan bool
}

func NewBindExe(closeChan chan bool) *BindExecutor {
	return nil
}
func (be *BindExecutor) Start() {

}
func (be *BindExecutor) Close() {

}
func (be *BindExecutor) MainRoutine() {

}
