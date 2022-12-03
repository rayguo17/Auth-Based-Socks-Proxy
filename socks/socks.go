package socks

import (
	"errors"
	"fmt"
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

var SupportedAuthMethod []byte = []byte{
	0, 2,
}

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
		fmt.Println("building client command")
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
func CheckAuthMethod(method byte) bool {
	found := false
	for i := 0; i < len(SupportedAuthMethod); i++ {
		if method == SupportedAuthMethod[i] {
			found = true
		}
	}
	return found

}
