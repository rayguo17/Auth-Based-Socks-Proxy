package connection

import (
	"net"
)

type SocksConnection struct {
	conn        net.Conn
	cmdExecutor Cmd
}
