package main

import (
	"context"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	Auth_None     int = 1
	Auth_Username     = 2
)

func main() {
	fmt.Println("Hello world")

	listener, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return
	}
	// do some initialization
	//main routine do something.
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting", err.Error())
			return
		}
		go acceptHandler(conn)

	}
}

var SupportedAuthMethod []byte = []byte{
	0, 1,
}

func acceptHandler(conn net.Conn) {
	buf := make([]byte, 512)
	initLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("err read", err.Error())
		return
	}
	fmt.Printf("%d bytes received!\n", initLen)
	count := 0
	authMethod := make([]int, 0)
	for i := 0; i < initLen; i++ {
		if i == 0 && buf[i] != 5 {
			fmt.Println("unsportted version! ending")
			return
		}
		if i == 1 {
			count = int(buf[i]) //how many method supported
		}
		if i > 1 {
			//check supported method
			if !checkAuthMethod(buf[i]) {
				fmt.Println("auth method not supported,ending...")
				return
			}
			//supported auth method. choose the bigger one
			authMethod = append(authMethod, int(buf[i]))

		}
	}
	//loop over choose which one to support
	chosenAuthMethod := 0
	for i := 0; i < len(authMethod); i++ {
		if authMethod[i] > chosenAuthMethod {
			chosenAuthMethod = authMethod[i]
		}
	}

	//2. then construct the response to the message bytes
	/*
		choose from supported method (later)
		directly support 0

	*/
	switch chosenAuthMethod {
	case 0:
		_, err = conn.Write([]byte{5, 1, 0})
		if err != nil {
			return
		}
	case 1:
		_, err = conn.Write([]byte{5, 1, 1})
		if err != nil {
			return
		}
	}

	//3. authentication phase (skip for now)
	if chosenAuthMethod == 1 {
		handleAuth(conn)
	}
	//4. request phase
	cmdBuf := make([]byte, 512)
	cmdLen, err := conn.Read(cmdBuf)
	if err != nil {
		return
	}

	address := strings.Builder{}

	//only supporting bind command for now.
	for i := 0; i < cmdLen; i++ {
		//fmt.Printf("%v ", cmdBuf[i])
		if i == 4 {
			count = int(cmdBuf[i])
			address.Write(cmdBuf[i+1 : i+1+count])

		}

	}
	portByte := cmdBuf[cmdLen-2 : cmdLen]
	pp.Println(portByte)
	portNum := int(portByte[0])*16*16 + int(portByte[1])
	fmt.Println("portNum", portNum)
	fmt.Println("")
	fmt.Println(address.String())
	fmt.Println("")
	//ver,cmd,rsv,ATYP,DST.ADDR                                            ,DST.PORT
	//  5,  1,  0,   3,      12,119,119,119,46,98,105,110,103,46,99,111,109, 1,187
	dialConn, err := net.Dial("tcp", address.String()+":"+strconv.Itoa(portNum))
	if err != nil {
		//send failure message back to client.
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer dialConn.Close()

	//send command and then start streaming??
	//return value:
	//ver, rep, RSV, ATYP, BND.ADDR, BND.PORT
	// 05,  00,  00,   01, 1, 0,0,0
	conn.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0})
	//use a go routine to send
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_, _ = io.Copy(conn, dialConn)
		//calculate
		cancel()
	}()
	go func() {
		_, _ = io.Copy(dialConn, conn)
		cancel()
	}()
	<-ctx.Done()
	return
}
func handleAuth(conn net.Conn) {

}
func checkAuthMethod(method byte) bool {
	found := false
	for i := 0; i < len(SupportedAuthMethod); i++ {
		if method == SupportedAuthMethod[i] {
			found = true
		}
	}
	return found

}

//browser tends to send multiple tcp connections, so there will be parallel thread using same instance.
