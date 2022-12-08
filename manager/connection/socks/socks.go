package socks

import (
	"net"
)

type SocksConnection struct {
	conn        net.Conn
	cmdExecutor Cmd
}
