package light

import (
	"errors"
	"github.com/rayguo17/go-socks/manager"
	"github.com/rayguo17/go-socks/manager/connection"
	"github.com/rayguo17/go-socks/util"
	"github.com/rayguo17/go-socks/util/protocol/light"
	"net"
	"time"
)

// 1.req auth to username:password
//version ulen uname plen pname
//  1       1    v     1    v
// 2.res auth
//version status
//   1      1
//          0/1 success/fail

func Authentication(ar *light.AuthReq, conn net.Conn) (*connection.AcpCon, error) {
	id := conn.RemoteAddr().String()
	username := string(ar.Uname)
	passwd := string(ar.Passwd)
	comm := manager.UM.GetConCommunicator()
	//fmt.Println(username, passwd)
	acpCon := connection.NewCon(id, conn, username, passwd, comm)
	manager.UM.AddCon(&acpCon)
	select {
	case authStatus := <-acpCon.AuthChan:
		if !authStatus {
			return nil, errors.New("authentication failed auth info incorrect")
		}
	case <-time.After(time.Second * 5):
		return nil, errors.New("Authenticate manager timeout")
	}
	return &acpCon, nil
}

func ConnectHandle(cmd *light.Cmd, con *connection.AcpCon) error {
	var addr util.Address
	switch int(cmd.AddrType) {
	case light.IPV4Address:
		addr = util.NewIpv4Addr(cmd.DstAddr, cmd.DstPort)
	case light.Domain:
		addr = util.NewDomainAddr(cmd.DstAddr, cmd.DstPort)
	default:
		return errors.New("addr type not supported")
	}
	return con.ConnectCmd(addr)
}
