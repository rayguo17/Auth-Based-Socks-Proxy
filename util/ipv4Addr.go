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
func (ia *Ipv4Addr) Port() [2]byte {
	return ia.port
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
	position := 0
	for i := len(ia.port) - 1; i >= 0; i-- {
		num := int(ia.port[i])

		for j := 0; j < position; j++ {
			num = num * 16 * 16
		}
		position++
		portNum += num

	}
	sb.WriteString(strconv.Itoa(portNum))
	return sb.String()
}
