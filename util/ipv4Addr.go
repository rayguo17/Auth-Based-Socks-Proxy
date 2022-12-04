package util

import (
	"strconv"
	"strings"
)

type Ipv4Addr struct {
	addr []byte
	port [2]byte
}

func NewIpv4Addr(ip []byte, port [2]byte) *Ipv4Addr {
	return &Ipv4Addr{
		addr: ip,
		port: port,
	}
}
func (ia *Ipv4Addr) Addr() string {
	sb := strings.Builder{}
	for i := 0; i < len(ia.addr); i++ {
		sb.WriteString(strconv.Itoa(int(ia.addr[i])))
		if i < len(ia.addr)-1 {
			sb.WriteRune('.')
		}
	}
	return sb.String()
}
func (ia *Ipv4Addr) String() string {
	sb := strings.Builder{}
	for i := 0; i < len(ia.addr); i++ {
		sb.WriteString(strconv.Itoa(int(ia.addr[i])))
		if i < len(ia.addr)-1 {
			sb.WriteRune('.')
		}
	}
	sb.WriteRune(':')
	var portNum int
	for i := 0; i < len(ia.port); i++ {
		portNum += int(ia.port[i])
	}
	sb.WriteString(strconv.Itoa(portNum))
	return sb.String()
}
