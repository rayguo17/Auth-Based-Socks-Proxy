package user

import (
	"context"
	"fmt"
	"github.com/rayguo17/go-socks/util"
	"io"
	"log"
	"net"
)

const (
	RUNNING int = 1
	DEAD    int = 2
	END     int = 3
)

type ConnectExecutor struct {
	targetCon net.Conn
	readChan  chan []byte
	writeChan chan []byte

	acpCon           *AcpCon
	chanClosedParent chan bool //acpCon inform executor to close routine
	toCloseChan      chan bool
	status           int
	addr             util.Address
}

func NewConExe(closeChan chan bool, conn net.Conn, addr util.Address, con *AcpCon) *ConnectExecutor {
	return &ConnectExecutor{
		chanClosedParent: closeChan,
		readChan:         make(chan []byte),
		writeChan:        make(chan []byte),
		targetCon:        conn,
		status:           DEAD,
		addr:             addr,
		toCloseChan:      make(chan bool),
		acpCon:           con,
	}
}
func (ce *ConnectExecutor) FormByte() []byte {
	//proxy port should be able to configure
	buf := []byte{5, 0, 0, 1, 127, 0, 0, 1, 13, 88}
	return buf
}
func (ce *ConnectExecutor) Start() error {
	ce.status = RUNNING
	//pp.Println(ce)
	ctx, cancel := context.WithCancel(context.Background())
	go ce.uploadRoutine(cancel)
	go ce.downloadRoutine(cancel)
	go ce.MainRoutine(ctx)
	return nil
}
func (ce *ConnectExecutor) Close() {
	switch ce.status {
	case RUNNING:
		ce.targetCon.Close()
		fmt.Println("close executed")
		ce.toCloseChan <- true
		ce.chanClosedParent <- true
	case DEAD:
		ce.targetCon.Close()
		ce.chanClosedParent <- true
	case END:
		ce.targetCon.Close()
	}
}

//only calculate traffic at the end?
func (ce *ConnectExecutor) uploadRoutine(cancelFunc context.CancelFunc) {
	written, err := io.Copy(ce.acpCon.conn, ce.targetCon)
	if err != nil {
		log.Print("upload routine err ")
		log.Println(err)
	}
	fmt.Printf("%d bytes written\n", written)
	cancelFunc()
}
func (ce *ConnectExecutor) downloadRoutine(cancelFunc context.CancelFunc) {
	written, err := io.Copy(ce.targetCon, ce.acpCon.conn)
	if err != nil {
		log.Print("download routine err ")
		log.Println(err)
	}
	fmt.Printf("%d bytes written\n", written)
	cancelFunc()
}
func (ce *ConnectExecutor) MainRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ce.status = END
			ce.Close()
			ce.acpCon.EndCommand() //might lead to multiple free???
			log.Println("executor main routine done")
			return
		case <-ce.toCloseChan:
			return

		}
	}
}