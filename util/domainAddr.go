package util

import (
	"strconv"
	"strings"
)

type DomainAddr struct {
	addr []byte
	port [2]byte
}

func NewDomainAddr(ip []byte, port [2]byte) *DomainAddr {
	return &DomainAddr{
		addr: ip,
		port: port,
	}
}
func (da *DomainAddr) String() string {
	sb := strings.Builder{}
	sb.Write(da.addr)
	sb.WriteRune(':')
	portNum := 0
	for i := 0; i < len(da.port); i++ {
		portNum += int(da.port[i])
	}
	sb.WriteString(strconv.Itoa(portNum))
	return sb.String()
}
