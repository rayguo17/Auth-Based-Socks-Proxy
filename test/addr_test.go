package test

import (
	"fmt"
	"github.com/rayguo17/go-socks/util"
	"strings"
	"testing"
)

func TestAddr(t *testing.T) {
	addr := "127.0.0.1:1009"
	addrArr := strings.Split(addr, ":")
	ipv4, _ := util.Ipv4FromString(addrArr[0], addrArr[1])
	str := ipv4.String()
	fmt.Println(str)
}
func TestPort(t *testing.T) {
	port := "1009"
	ports, _ := util.ParsePort(port)
	fmt.Println(ports)
}
