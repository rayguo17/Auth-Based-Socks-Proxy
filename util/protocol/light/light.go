package light

import (
	"errors"
	"fmt"
	"github.com/rayguo17/go-socks/util"
)

const (
	IPV4Address = 1
	Domain      = 2
	IPV6Address = 3
)

type AuthReq struct {
	Ulen   uint8
	Uname  []byte
	Plen   uint8
	Passwd []byte
}

type Cmd struct {
	AddrType uint8
	DstAddr  []byte
	DstPort  [2]byte
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
		Passwd: buf[2+ulen:],
	}
	return res, nil
}
func BuildCmd(buf []byte) (*Cmd, error) {
	res := &Cmd{}
	switch buf[0] {
	case IPV4Address:
		if len(buf) != 7 {
			return nil, errors.New(fmt.Sprintf("client command length does not match address type: %s, expect: %d, got: %d", util.Ipv4Address, 7, len(buf)))
		}
		res.AddrType = buf[0]
		res.DstAddr = make([]byte, 0, 4)
		for i := 0; i < 4; i++ {
			res.DstAddr = append(res.DstAddr, buf[1+i])
		}
	case Domain:
		addrLen := int(buf[1])
		if len(buf) != addrLen+4 {
			return nil, errors.New(fmt.Sprintf("client command length does not match address type: %s, expect: %d, got: %d", util.DomainAddress, 4+addrLen, len(buf)))
		}
		res.AddrType = buf[0]
		res.DstAddr = make([]byte, 0, addrLen)
		for i := 0; i < addrLen; i++ {
			res.DstAddr = append(res.DstAddr, buf[2+i])
		}
	case IPV6Address:
		if len(buf) != 19 {
			return nil, errors.New(fmt.Sprintf("client command length does not match address type: %s, expect: %d, got: %d", util.Ipv6Address, 19, len(buf)))
		}
		res.AddrType = buf[0]
		res.DstAddr = make([]byte, 0, 16)
		for i := 0; i < 16; i++ {
			res.DstAddr = append(res.DstAddr)
		}
	default:
		return nil, errors.New("address type unrecognized")
	}
	bufLen := len(buf)
	res.DstPort = [2]byte{buf[bufLen-2], buf[bufLen-1]}
	return res, nil
}

func FormAR(uname string, passwd string) []byte {
	nameArr := []byte(uname)
	passArr := []byte(passwd)
	res := make([]byte, 0, 2+len(nameArr)+len(passArr))
	res = append(res, byte(len(nameArr)))
	for i := 0; i < len(nameArr); i++ {
		res = append(res, nameArr[i])
	}
	res = append(res, byte(len(passArr)))
	for i := 0; i < len(passArr); i++ {
		res = append(res, passArr[i])
	}
	return res
}

func FormCmd(address util.Address) []byte {
	addrType := 0
	switch address.AddrType() {
	case util.Ipv4Address:
		addrType = IPV4Address
	case util.DomainAddress:
		addrType = Domain
	case util.Ipv6Address:
		addrType = IPV6Address
	default:
		addrType = 0
	}
	addrByte := address.AddrByte()
	portByte := address.Port()
	res := make([]byte, 0, 3+len(addrByte))

	res = append(res, byte(addrType))
	if addrType == Domain {
		res = append(res, byte(len(addrByte)))
	}
	for i := 0; i < len(addrByte); i++ {
		res = append(res, addrByte[i])
	}
	res = append(res, portByte[0], portByte[1])
	return res
}
