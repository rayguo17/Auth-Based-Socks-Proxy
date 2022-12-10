package socks

import (
	"errors"
	"fmt"
	"github.com/rayguo17/go-socks/manager"
	"github.com/rayguo17/go-socks/manager/connection"
	"github.com/rayguo17/go-socks/util"
	"github.com/rayguo17/go-socks/util/protocol/socks"
	"net"
	"time"
)

//how to encapsulate protocol??

//light response with the supported light

func CommandHandle(cmd *socks.ClientCmd, con *connection.AcpCon) error {
	switch cmd.Cmd {
	case socks.BIND:
		//fmt.Println("BIND command")

	case socks.CONNECT:
		//fmt.Println("CONNECT command")
		return handleConnect(cmd, con)
	case socks.UDPASSO:
		//fmt.Println("UDP ASSOCIATE command")
	}
	return nil
}

//according user rule, choose where to go.
func handleConnect(cmd *socks.ClientCmd, con *connection.AcpCon) error {

	var addr util.Address
	switch int(cmd.Atyp) {
	case socks.Ipv4Address:
		addr = util.NewIpv4Addr(cmd.DstAddr, cmd.DstPort)
	case socks.DomainAddress:
		//fmt.Println("domain address")
		addr = util.NewDomainAddr(cmd.DstAddr, cmd.DstPort)
	}
	//fmt.Println(addr.String())
	return con.ConnectCmd(addr)

}

func Authenticate(authReq *socks.AuthReq, conn net.Conn) (*connection.AcpCon, error) {
	id := conn.RemoteAddr().String()
	username := string(authReq.Uname)
	passwd := string(authReq.Passwd)
	comm := manager.UM.GetConCommunicator()
	acpCon := connection.NewCon(id, conn, username, passwd, comm)
	manager.UM.AddCon(&acpCon)
	select {
	case authStatus := <-acpCon.AuthChan:
		if !authStatus {
			return nil, errors.New("authentication failed")
		}
		//fmt.Println("acpCon testing", acpCon)
	case <-time.After(time.Second * 5):
		return nil, errors.New("Authenticate manager timeout")
	}

	return &acpCon, nil
}
func ConstructResp(con *connection.AcpCon, msgType int) ([]byte, error) {
	switch msgType {
	case socks.HandShakeResponse:
		fmt.Println("building handshake response")
	case socks.AuthRequest:
		fmt.Println("building Auth Response")
	case socks.ServerResponse:
		return con.CmdResponse()
	default:
		return nil, errors.New(fmt.Sprintf("unrecognized msgType: %d\n", msgType))
	}
	return nil, nil

}
func Anonymous(conn net.Conn) (*connection.AcpCon, error) {
	id := conn.RemoteAddr().String()
	username := "anonymous"
	passwd := ""
	comm := manager.UM.GetConCommunicator()
	acpCon := connection.NewCon(id, conn, username, passwd, comm)
	manager.UM.AddCon(&acpCon)
	authStatus := <-acpCon.AuthChan
	if !authStatus {
		return nil, errors.New("authentication failed")
	}
	return &acpCon, nil
}
