package light

import "errors"

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
