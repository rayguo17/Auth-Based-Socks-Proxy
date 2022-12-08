package light

import (
	"errors"
	"github.com/rayguo17/go-socks/manager/connection/socks"
)

// 1.req auth to username:password
//version ulen uname plen pname
//  1       1    v     1    v
// 2.res auth
//version status
//   1      1
//          0/1 success/fail

type AuthReq struct {
	Ulen   uint8
	Uname  []byte
	Plen   uint8
	Passwd []byte
}

func BuildAR(buf []byte) (*AuthReq, error) {
	//check if ulen and plen match
	size := len(buf)
	ulen := int(buf[0])
	plen := int(buf[1+ulen])
	if size != ulen+plen+2 {
		return nil, errors.New("ulen+plen not match with buf size")
	}
	res := &AuthReq{
		Ulen:   buf[0],
		Plen:   buf[1+ulen],
		Uname:  buf[1 : 1+ulen],
		Passwd: buf[1+ulen:],
	}
	return res, nil
}

func Authentication(ar *AuthReq) (*socks.AcpCon, error) {

	return nil, nil
}
