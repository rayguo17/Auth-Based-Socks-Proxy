package util

import (
	"strconv"
	"strings"
)

type DomainAddr struct {
	addr []byte
	port [2]byte
}

func (da *DomainAddr) AddrType() string {
	return DomainAddress
}

func NewDomainAddr(ip []byte, port [2]byte) *DomainAddr {
	return &DomainAddr{
		addr: ip,
		port: port,
	}
}
func (da *DomainAddr) Addr() string {
	sb := strings.Builder{}
	sb.Write(da.addr)
	return sb.String()
}
func (da *DomainAddr) AddrByte() []byte {
	return da.addr
}
func (da *DomainAddr) Port() [2]byte {
	return da.port
}
func (da *DomainAddr) String() string {
	sb := strings.Builder{}
	sb.Write(da.addr)
	sb.WriteRune(':')
	//pp.Println(da.port)
	portNum := 0
	position := 0
	for i := len(da.port) - 1; i >= 0; i-- {
		num := int(da.port[i])

		for j := 0; j < position; j++ {
			num = num * 16 * 16
		}
		position++
		portNum += num

	}
	sb.WriteString(strconv.Itoa(portNum))
	return sb.String()
}
