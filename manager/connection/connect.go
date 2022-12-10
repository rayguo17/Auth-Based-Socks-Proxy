package connection

import (
	"context"
	"fmt"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/manager/common"
	"github.com/rayguo17/go-socks/util"
	"io"
	"log"
	"net"
	"time"
)

//how to organize program. how to abstract code.
const (
	RUNNING int = 1
	DEAD    int = 2
	END     int = 3
)

var EXECSTATUS = map[int]string{
	0: "NULL",
	1: "RUNNING",
	2: "DEAD",
	3: "END",
}

type ConnectExecutor struct {
	targetCon        net.Conn
	readChan         chan []byte
	writeChan        chan []byte
	acpCon           *AcpCon
	isRemote         bool
	RemoteAddr       util.Address
	chanClosedParent chan bool //acpCon inform executor to close routine
	toCloseChan      chan bool
	status           int
	DstAddr          util.Address
}

func (ce *ConnectExecutor) Status() int {
	return ce.status //should put it inside routine
}
func (ce *ConnectExecutor) RemoteAddress() string {
	return ce.DstAddr.String()
}

func NewConExe(closeChan chan bool, isRemote bool, RemoteAddr util.Address, conn net.Conn, DstAddr util.Address, con *AcpCon) *ConnectExecutor {
	return &ConnectExecutor{
		chanClosedParent: closeChan,
		readChan:         make(chan []byte),
		writeChan:        make(chan []byte),
		targetCon:        conn,
		isRemote:         isRemote,
		RemoteAddr:       RemoteAddr,
		status:           DEAD,
		DstAddr:          DstAddr,
		toCloseChan:      make(chan bool),
		acpCon:           con,
	}
}
func (ce *ConnectExecutor) FormByte() []byte {
	//proxy port should be able to configure
	portNum := config.SystemConfig.SocksPort
	portByte := util.DecToByte(portNum)
	//pp.Println(byteArr)
	buf := []byte{5, 0, 0, 1, 127, 0, 0, 1, portByte[0], portByte[1]} //port number
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
	written, err := io.Copy(ce.acpCon.GetConn(), ce.targetCon)
	if err != nil {
		//log.Print("upload routine err ")
		//log.Println(err)
	}
	inform := ce.uploadTraffic(written, true)
	select {
	case _ = <-inform:
		//log.Println(res.String())
	case <-time.After(time.Second * 5):
		log.Println("upload traffic timeout quiting anyway")
	}
	//fmt.Printf("%d bytes written\n", written)
	cancelFunc()
}
func (ce *ConnectExecutor) downloadRoutine(cancelFunc context.CancelFunc) {
	written, err := io.Copy(ce.targetCon, ce.acpCon.GetConn())
	if err != nil {
		//log.Print("download routine err ")
		//log.Println(err)
	}
	inform := ce.uploadTraffic(written, false)
	select {
	case _ = <-inform:
		//log.Println(res.String())
	case <-time.After(time.Second * 5):
		log.Println("upload traffic timeout quiting anyway")
	}
	//fmt.Printf("%d bytes written\n", written)
	cancelFunc()
}
func (ce *ConnectExecutor) uploadTraffic(traffic int64, upload bool) chan *util.Response {
	res := make(chan *util.Response)

	wrap := &common.UploadTrafficWrap{
		Username:   ce.acpCon.GetName(),
		Up:         upload,
		Count:      traffic,
		InformChan: res,
	}
	ce.acpCon.UploadTraffic(wrap)
	return res
}
func (ce *ConnectExecutor) MainRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ce.status = END
			ce.Close()
			ce.acpCon.EndCommand() //might lead to multiple free???
			//log.Println("executor main routine done")
			return
		case <-ce.toCloseChan:
			return

		}
	}
}
