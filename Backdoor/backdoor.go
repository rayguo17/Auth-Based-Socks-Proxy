package Backdoor

import (
	"bufio"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"github.com/rayguo17/go-socks/manager"
	"os"
)

func BackDoorRoutine() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		str, _ := inputReader.ReadString('\n')
		//fmt.Println(str)
		switch str {
		case "1\n":
			manager.UM.ListUsers()
		case "2\n":
			manager.UM.ListConn()
			//pp.Println(manager.UM.AcpConnections)
		case "q\n":
			return
		case "3\n":
			pp.Println(manager.UM.Users)
		default:
			fmt.Println("cmd unrecognized")
		}
	}
}
