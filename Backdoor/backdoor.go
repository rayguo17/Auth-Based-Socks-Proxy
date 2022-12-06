package Backdoor

import (
	"bufio"
	"fmt"
	"github.com/rayguo17/go-socks/user"
	"os"
)

func BackDoorRoutine() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		str, _ := inputReader.ReadString('\n')
		fmt.Println(str)
		switch str {
		case "1\n":
			user.UM.ListUsers()
		case "2\n":
			user.UM.ListConn()
			//pp.Println(user.UM.AcpConnections)
		default:
			return
		}
	}
}
