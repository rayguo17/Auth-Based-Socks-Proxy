/*
	Package Connections
	We store each tcp listener accepted connection as connections
*/
package connections

import "github.com/rayguo17/go-socks/user"

type AcpCon struct {
	id         string // identifier ("address:port")
	owner      *user.User
	bytesCount int
	AuthChan   chan bool
}

func NewCon() {

}
