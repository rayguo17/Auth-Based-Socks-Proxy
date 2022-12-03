package Backdoor

import (
	"bufio"
	"fmt"
	"github.com/k0kubun/pp/v3"
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
			user.UM.ListUser()
			fmt.Println("print")
		case "2\n":
			pp.Println(user.UM.AcpConnections)
		default:
			return
		}
	}
}
