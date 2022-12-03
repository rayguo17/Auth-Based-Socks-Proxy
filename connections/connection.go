/*
	Package Connections
	We store each tcp listener accepted connection as connections
*/
package connections

import "github.com/rayguo17/go-socks/user"

type AcpCon struct {
	owner      *user.User
	bytesCount int
}
