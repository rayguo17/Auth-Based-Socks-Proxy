package util

import (
	"strconv"
	"strings"
)

type Ipv4Addr struct {
	addr []byte
	port [2]byte
}

func Ipv4FromString(ip string, port string) (*Ipv4Addr, error) {
	portNums, err := ParsePort(port)
	if err != nil {
		return nil, err
	}
	ipArr := strings.Split(ip, ".")
	ipByte := make([]byte, 0, 4)
	for i := 0; i < len(ipArr); i++ {
		num, err := strconv.Atoi(ipArr[i])
		if err != nil {
			return nil, err
		}
		ipByte = append(ipByte, byte(num))
	}
	return &Ipv4Addr{
		addr: ipByte,
		port: portNums,
	}, nil
}
func (ia *Ipv4Addr) AddrType() string {
	return Ipv4Address
}
func DecToByte(port int) [2]byte {

	first := port / 256
	second := port % 256
	return [2]byte{byte(first), byte(second)}

}
func ParsePort(port string) ([2]byte, error) {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return [2]byte{}, err
	}
	return DecToByte(portNum), nil
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
func (ia *Ipv4Addr) AddrByte() []byte {
	return ia.addr
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
