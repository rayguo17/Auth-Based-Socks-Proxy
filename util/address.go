package util

import "strings"

type Address interface {
	String() string
}

type Addr struct {
	addr []byte
	port [2]byte
}

func NewAddr(addr []byte, port [2]byte) *Ipv4Addr {
	return &Ipv4Addr{
		addr: addr,
		port: port,
	}
}
func (a *Addr) String() string {
	sb := strings.Builder{}
	sb.Write(a.addr)
	sb.WriteRune(':')
	for i := 0; i < len(a.port); i++ {
		sb.WriteByte(a.port[i])
	}
	return sb.String()
}
