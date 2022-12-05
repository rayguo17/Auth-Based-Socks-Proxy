package socks

import (
	"errors"
	"fmt"
	"github.com/rayguo17/go-socks/user"
	"github.com/rayguo17/go-socks/util"
	"net"
)

//how to encapsulate protocol??
const (
	HandShakeRequest  int = 1
	HandShakeResponse     = 2
	AuthRequest           = 3
	AuthResponse          = 4
	ClientCommand         = 5
	ServerResponse        = 6
)

//Manager Address through interface??
const (
	Ipv4Address   int = 1
	DomainAddress     = 3
	Ipv6Address       = 4
)

//Manage command through intercace?
const (
	CONNECT uint8 = 1
	BIND    uint8 = 2
	UDPASSO uint8 = 3
)

var SupportedAuthMethod []byte = []byte{
	0, 2,
}
var SupportedCmd []byte = []byte{1}

type HandshakeReq struct {
	Version    uint8
	AuthCount  uint8
	AuthMethod []byte
}

//server response with the supported server
type HandShakeResp struct {
	Version    uint8
	AuthMethod uint8
}
type AuthReq struct {
	Version uint8
	Ulen    uint8
	Uname   []byte
	Plen    uint8
	Passwd  []byte
}
type AuthResp struct {
	Version uint8
	Status  uint8
}
type ClientCmd struct {
	Version uint8
	Cmd     uint8
	Rsv     uint8 //default to 0
	Atyp    uint8
	DstAddr []byte //depends on atyp
	DstPort [2]byte
}
type ServerResp struct {
	Version uint8
	Rep     uint8 //status code
	Rsv     uint8
	Atyp    uint8
	BndAddr []byte
	BndPort [2]byte
}

func FromByte(buf []byte, msgType int) (interface{}, error) {
	//
	switch msgType {
	case HandShakeRequest:
		return buildHSR(buf)
	case HandShakeResponse:
		fmt.Println("building handshake response")
	case AuthRequest:
		return buildAR(buf)
	case AuthResponse:
		fmt.Println("building Auth Response")
	case ClientCommand:
		return buildCC(buf)
	case ServerResponse:
		fmt.Println("building server response")
	default:
		return nil, errors.New(fmt.Sprintf("unrecognized msgType: %d\n", msgType))
	}
	return nil, nil
}
func FromStruct(target interface{}, msgType int) ([]byte, error) {
	switch v := target.(type) {
	case HandshakeReq:
		fmt.Println(v)
	case HandShakeResp:
		fmt.Println(v)
	case AuthReq:
		fmt.Println(v)
	case AuthResp:
		fmt.Println(v)
	case ClientCmd:
		fmt.Println(v)
	case ServerResp:
		fmt.Println()
	default:
		fmt.Println("unrecognized type interface,quiting!!")

	}
	return nil, nil
}
func CommandHandle(cmd *ClientCmd, con *user.AcpCon) error {
	switch cmd.Cmd {
	case BIND:
		//fmt.Println("BIND command")

	case CONNECT:
		//fmt.Println("CONNECT command")
		return handleConnect(cmd, con)
	case UDPASSO:
		//fmt.Println("UDP ASSOCIATE command")
	}
	return nil
}
func handleConnect(cmd *ClientCmd, con *user.AcpCon) error {
	var addr util.Address
	switch int(cmd.Atyp) {
	case Ipv4Address:
		addr = util.NewIpv4Addr(cmd.DstAddr, cmd.DstPort)
	case DomainAddress:
		//fmt.Println("domain address")
		addr = util.NewDomainAddr(cmd.DstAddr, cmd.DstPort)
	}
	//fmt.Println(addr.String())
	return con.ConnectCmd(addr)

}
func Authenticate(authReq *AuthReq, conn net.Conn) (*user.AcpCon, error) {
	id := conn.RemoteAddr().String()
	username := string(authReq.Uname)
	passwd := string(authReq.Passwd)
	acpCon := user.NewCon(id, conn, username, passwd)
	user.UM.AddCon(&acpCon)
	authStatus := <-acpCon.AuthChan
	if !authStatus {
		return nil, errors.New("authentication failed")
	}
	//fmt.Println("acpCon testing", acpCon)

	return &acpCon, nil
}
func ConstructResp(con *user.AcpCon, msgType int) ([]byte, error) {
	switch msgType {
	case HandShakeResponse:
		fmt.Println("building handshake response")
	case AuthResponse:
		fmt.Println("building Auth Response")
	case ServerResponse:
		return con.CmdResponse()
	default:
		return nil, errors.New(fmt.Sprintf("unrecognized msgType: %d\n", msgType))
	}
	return nil, nil

}
func Anonymous(conn net.Conn) (*user.AcpCon, error) {
	id := conn.RemoteAddr().String()
	username := "anonymous"
	passwd := ""
	acpCon := user.NewCon(id, conn, username, passwd)
	user.UM.AddCon(&acpCon)
	authStatus := <-acpCon.AuthChan
	if !authStatus {
		return nil, errors.New("authentication failed")
	}
	return &acpCon, nil
}

func buildCC(buf []byte) (*ClientCmd, error) {
	//pp.Println(buf)
	res := &ClientCmd{}
	if buf[0] != 5 {
		return nil, errors.New(fmt.Sprintf("socks version: %d not supported\n", buf[0]))
	}
	res.Version = buf[0]
	if !checkCmd(buf[1]) {
		return nil, errors.New(fmt.Sprintf("socks command: %d not supported\n", buf[1]))
	}
	res.Cmd = buf[1]
	switch int(buf[3]) {
	case Ipv4Address:
		if len(buf) != 4+4+2 {
			return nil, errors.New(fmt.Sprintf("client command length does not match address type: %d, expect: %d, got: %d", Ipv4Address, 10, len(buf)))
		}
		res.Atyp = buf[3]
		res.DstAddr = make([]byte, 0, 4)
		for i := 0; i < 4; i++ {
			res.DstAddr = append(res.DstAddr, buf[4+i])
		}
	case DomainAddress:
		addressLen := int(buf[4])
		if len(buf) != 5+addressLen+2 {
			return nil, errors.New(fmt.Sprintf("client command length does not match address type: %d, expect: %d, got: %d", DomainAddress, 7+addressLen, len(buf)))
		}
		res.Atyp = buf[3]

		res.DstAddr = make([]byte, 0, addressLen)
		for i := 0; i < addressLen; i++ {
			res.DstAddr = append(res.DstAddr, buf[5+i])
		}
	case Ipv6Address:
		if len(buf) != 4+16+2 {
			return nil, errors.New(fmt.Sprintf("client command length does not match address type: %d, expect: %d, got: %d", Ipv6Address, 22, len(buf)))
		}
		res.Atyp = buf[3]
		res.DstAddr = make([]byte, 0, 16)
		for i := 0; i < 16; i++ {
			res.DstAddr = append(res.DstAddr)
		}
	}
	bufLen := len(buf)
	res.DstPort = [2]byte{buf[bufLen-2], buf[bufLen-1]}
	return res, nil
}
func buildAR(buf []byte) (*AuthReq, error) {
	//interact with user
	res := &AuthReq{}
	if buf[0] != 1 {
		return nil, errors.New(fmt.Sprintf("auth version: %d not supported\n", buf[0]))
	}
	res.Version = buf[0]
	res.Ulen = buf[1]
	if res.Ulen > 255 {
		return nil, errors.New(fmt.Sprintf("socks ulen %d too big\n", buf[1]))
	}
	uname := make([]byte, 0, res.Ulen)
	ulen := int(res.Ulen)
	for i := 0; i < ulen; i++ {
		uname = append(uname, buf[i+2])
	}
	res.Uname = uname
	res.Plen = buf[2+ulen]
	plen := int(res.Plen)
	if plen != len(buf)-ulen-3 {
		return nil, errors.New(fmt.Sprintf("socks plen %d not matched with actual len %d\n", buf[1], len(buf)-ulen-3))
	}
	for i := 0; i < plen; i++ {
		res.Passwd = append(res.Passwd, buf[3+ulen+i])
	}
	return res, nil
}
func buildHSR(buf []byte) (*HandshakeReq, error) {
	res := &HandshakeReq{}
	if buf[0] != 5 {
		return nil, errors.New(fmt.Sprintf("socks version: %d not supported\n", buf[0]))
	}
	res.Version = buf[0]
	res.AuthCount = buf[1]
	if int(res.AuthCount) != len(buf)-2 {
		return nil, errors.New(fmt.Sprintf("socks auth count (%d) not the same as actual supported auth method (%d)",
			res.AuthCount, len(buf)-2))
	}
	for i := 2; i < len(buf); i++ {
		if !CheckAuthMethod(buf[i]) {
			continue
		}
		res.AuthMethod = append(res.AuthMethod, buf[i])

	}
	res.AuthCount = uint8(len(res.AuthMethod))
	return res, nil
}
func checkCmd(cmd byte) bool {
	for i := 0; i < len(SupportedCmd); i++ {
		if SupportedCmd[i] == cmd {
			return true
		}
	}
	return false
}
func CheckAuthMethod(method byte) bool {
	found := false
	for i := 0; i < len(SupportedAuthMethod); i++ {
		if method == SupportedAuthMethod[i] {
			found = true
		}
	}
	return found

}
