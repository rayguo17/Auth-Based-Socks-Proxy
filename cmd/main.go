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

func main() {
	fmt.Println("Hello world")

	listener, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting", err.Error())
			return
		}
		go acceptHandler(conn)

	}
}

var Supported_Auth_Method []byte = []byte{
	0, 1,
}

func acceptHandler(conn net.Conn) {
	//first initiate socks handshake
	//1. receive auth message:
	buf := make([]byte, 512)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Println("err read", err.Error())
		return
	}
	fmt.Printf("%d bytes received!\n", len)
	//pp.Println(buf)
	//parse the buf by iterate
	count := 0
	//support how many authenticate method??? maybe for now only user and none??
	// add user authentication later, we use config file to store user information.
	//used_auth_method := make([]byte,len)
	for i, b := range buf {
		if i == 0 && b != 5 {
			fmt.Println("unsportted version! ending")
			return
			//should be a go to if later we seperate authentication and accept handling
		}

		//fmt.Printf("index: %v, value: %v\n", i, b)
		count++
		if count == len {
			break
		}
	}

	//2. then construct the response to the message bytes
	/*
		choose from supported method (later)
		directly support 0

	*/
	_, err = conn.Write([]byte{5, 0})
	if err != nil {
		return
	}
	//3. authentication phase (skip for now)

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

//browser tends to send multiple tcp connections, so there will be parallel thread using same instance.
