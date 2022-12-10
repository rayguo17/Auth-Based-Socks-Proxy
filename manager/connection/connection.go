/*
	Package Connections
	We store each tcp listener accepted connection as connections
*/
package connection

import (
	"errors"
	"fmt"
	"github.com/rayguo17/go-socks/manager/common"
	"github.com/rayguo17/go-socks/util"
	"github.com/rayguo17/go-socks/util/logger"
	"log"
	"net"
	"time"
)

const (
	AuthDone int = 1
	CmdRecv      = 2
	Working      = 3
	End          = 4
	Dead         = 5
)

var ACPSTATUSMAP = map[int]string{
	1: "AuthDone",
	2: "CmdRecv",
	3: "Working",
	4: "End",
	5: "Dead",
}

//COMMAND TYPE

const (
	Connect int = 1
	Bind        = 2
	UDPAsso     = 3
)
const (
	Direct int = 1
	Light      = 2
)

var CMDMap = map[int]string{
	1: "Connect",
	2: "Bind",
	3: "UDPAsso",
}

const BUFMAX = 4096

type AcpCon struct {
	id            string // identifier ("address:port")
	bytesCount    int
	AuthChan      chan bool
	conn          net.Conn
	username      string //should be abstract to auth
	passwd        string
	cmdType       int
	cmdClosedChan chan bool //tell cmd to close
	isRemote      bool
	remoteAddr    util.Address
	cmdExecutor   Cmd
	status        int
	//we can store UM's channel to communicate with it.
	communicate *common.Communicator
}

func (acpCon *AcpCon) SetRemote(isRemote bool, remoteAddr util.Address) {
	acpCon.isRemote = isRemote
	acpCon.remoteAddr = remoteAddr
}

func NewCon(id string, conn net.Conn, username string, passwd string, comm *common.Communicator) AcpCon {
	return AcpCon{
		id:         id,
		AuthChan:   make(chan bool),
		conn:       conn,
		bytesCount: 0,
		//should have a auth interface...
		username:    username,
		passwd:      passwd,
		status:      AuthDone,
		communicate: comm,
	}
}

func (acpCon *AcpCon) Log() string {
	str := fmt.Sprintf("%v %v %v", acpCon.username, acpCon.id, acpCon.RemoteAddress())
	return str
}
func (acpCon *AcpCon) EndCommand() {
	if acpCon.status == Dead {
		//fmt.Println("already dead")
		return
	}
	//fmt.Println("Still alive ending")
	acpCon.handleEnd()
}
func (acpCon *AcpCon) CloseCommand() {
	if acpCon.status == Dead || acpCon.status == End {
		return
	}
	acpCon.handleClose()
}
func (acp *AcpCon) ProtocolClose() {
	if acp.status == Working || acp.status == End || acp.status == Dead {
		return
	}
	acp.handleClose()
}

//could be manually killed or by closing the socket.

//detail routine should be maintain by sub type
//delete main routine.... decouple with sub cmd.

func (acpCon *AcpCon) handleEnd() {
	acpCon.status = End
	//fmt.Println("handling end in acpCon")
	informChan := make(chan *util.Response)
	req := &common.DCWrap{
		Id:         acpCon.username + "|" + acpCon.id,
		InformChan: informChan,
	}
	acpCon.communicate.DelCon(req)
	//should add a timeout don't wait forever
	select {
	case <-informChan:
		//log.Println("acp con delete success")
	case <-time.After(5 * time.Second):
		log.Println("delete time out quiting anyway")
	}
	return
}

func (acpCon *AcpCon) handleClose() {
	if acpCon.status == Dead {
		return
	}
	acpCon.status = Dead
	//TODO:
	//end sub routine (if exist), delete from manager manager.
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
	informChan := make(chan *util.Response)
	req := &common.DCWrap{
		Id:         acpCon.username + "|" + acpCon.id,
		InformChan: informChan,
	}
	acpCon.communicate.DelCon(req)
	//should add a timeout don't wait forever
	select {
	case <-informChan:
		//log.Println("acp con delete success")
	case <-time.After(5 * time.Second):
		log.Println("delete time out quiting anyway")
	}
	return
}

//independent
func (acpCon *AcpCon) ConnectCmd(addr util.Address) error {
	acpCon.cmdType = Connect
	acpCon.status = CmdRecv
	addStr := addr.Addr()
	informChan := make(chan *util.Response)
	req := &common.CheckRulesetWrap{
		Username:   acpCon.username,
		DstAddr:    addStr,
		InformChan: informChan,
	}
	acpCon.communicate.CheckRuleset(req)
	select {
	case resp := <-informChan:
		if resp.GetErrCode() != 0 {
			logger.Access.Println(acpCon.Log() + " " + addStr + " " + resp.GetErrMsg())
			return errors.New(resp.GetErrMsg())
		}
	case <-time.After(time.Second * 5):
		return errors.New("check ruleset timeout")
	}
	str := addr.String()
	conn, err := net.DialTimeout("tcp", str, time.Second*10)
	if err != nil {
		return err
	}
	//success, create executor
	acpCon.cmdClosedChan = make(chan bool)

	connectExe := NewConExe(acpCon.cmdClosedChan, conn, addr, acpCon)
	acpCon.cmdExecutor = connectExe
	logger.Access.Println(acpCon.Log() + " accepted")
	//fmt.Println("command execute")
	return nil
}
func (acpCon *AcpCon) ExecuteBegin() error {
	//should check everything before begin
	err := acpCon.cmdExecutor.Start()
	if err != nil {
		return err
	}
	acpCon.status = Working
	return err
}
func (acpCon *AcpCon) UploadTraffic(wrap *common.UploadTrafficWrap) {
	acpCon.communicate.UploadTrrafic(wrap)
}
func (acpCon *AcpCon) CmdResponse() ([]byte, error) {
	if acpCon.cmdExecutor == nil {
		return nil, errors.New("command Executor has not been initialize")
	}
	return acpCon.cmdExecutor.FormByte(), nil
}
func (acp *AcpCon) RemoteAddress() string {
	if acp.cmdExecutor == nil {
		return ""
	} else {
		return acp.cmdExecutor.RemoteAddress()
	}
}
func (acp *AcpCon) ExecutorStatus() int {
	if acp.cmdExecutor == nil {
		return 0
	} else {
		return acp.cmdExecutor.Status()
	}
}
func (acp *AcpCon) GetConn() net.Conn {
	return acp.conn
}
func (acp *AcpCon) GetName() string {
	return acp.username
}
func (acp *AcpCon) GetStatus() int {
	return acp.status
}
func (acp *AcpCon) GetCmdType() int {
	return acp.cmdType
}
func (acp *AcpCon) GetPasswd() string {
	return acp.passwd
}
func (acp *AcpCon) GetID() string {
	return acp.id
}
