package main

import (
	"context"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"github.com/rayguo17/go-socks/Backdoor"
	"github.com/rayguo17/go-socks/socks"
	"github.com/rayguo17/go-socks/user"
	"io"
	"log"
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
	// do some initialization
	//main routine do something.
	go user.UM.MainRoutine()
	go Backdoor.BackDoorRoutine()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting", err.Error())
			return
		}
		go acceptHandler(conn)

	}
}

func acceptHandler(conn net.Conn) {
	buf := make([]byte, 512)
	initLen, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%d bytes received!\n", initLen)
	count := 0

	source, err := socks.FromByte(buf[:initLen], socks.HandShakeRequest)

	if err != nil {
		log.Println(err)
		return
	}
	if _, ok := source.(*socks.HandshakeReq); !ok {
		log.Println("socks handshake mapping failed returning")
		return
	}
	handShakeReq := source.(*socks.HandshakeReq)
	//loop over choose which one to support
	chosenAuthMethod := 0
	for i := 0; i < len(handShakeReq.AuthMethod); i++ {
		if int(handShakeReq.AuthMethod[i]) > chosenAuthMethod {
			chosenAuthMethod = int(handShakeReq.AuthMethod[i])
		}
	}
	//2. then construct the response to the message bytes
	/*
		choose from supported method (later)
		directly support 0

	*/
	//pp.Println(handShakeReq)
	//log.Printf("chosen auth method: %d", chosenAuthMethod)
	switch chosenAuthMethod {
	case 0:
		_, err = conn.Write([]byte{5, 0})
		if err != nil {
			log.Println(err)
			return
		}
	case 2:
		_, err = conn.Write([]byte{5, 2})
		if err != nil {
			log.Println(err)
			return
		}
	}
	//3. authentication phase (skip for now)
	var acpCon *user.AcpCon
	if chosenAuthMethod == 2 {
		authBuf := make([]byte, 512)
		authLen, err := conn.Read(authBuf)
		if err != nil {
			log.Println(err)
			return
		}
		//register tmp to received feedback.

		source, err = socks.FromByte(authBuf[:authLen], socks.AuthRequest)
		if err != nil {
			log.Println(err)
			return
		}
		if _, ok := source.(*socks.AuthReq); !ok {
			log.Fatal("socks authReq mapping failed returning")
			return
		}
		authReq := source.(*socks.AuthReq)
		acpCon, err = socks.Authenticate(authReq, conn)
		if err != nil {
			log.Println(err)
			conn.Write([]byte{1, 1})
			return
		}
	} else {
		//create anonymous account route
		acpCon, err = socks.Anonymous(conn)
		if err != nil {
			log.Println(err)
			conn.Write([]byte{1, 1})
			return
		}
	}
	defer acpCon.ManualClose() //close on return
	_, err = conn.Write([]byte{1, 0})
	if err != nil {
		log.Println(err)
		return
	}
	//TODO://4. request phase
	cmdBuf := make([]byte, 512)
	cmdLen, err := conn.Read(cmdBuf)
	if err != nil {
		log.Println(err)

		return
	}
	source, err = socks.FromByte(cmdBuf[:cmdLen], socks.ClientCommand)
	if err != nil {
		log.Println(err)
		return
	}
	if _, ok := source.(*socks.ClientCmd); !ok {
		log.Fatal("socks authReq mapping failed returning")
		return
	}
	clientCmd := source.(*socks.ClientCmd)
	//pp.Println(clientCmd)
	//commandHandle -> handleConnect -> con.ConnectCmd
	err = socks.CommandHandle(clientCmd, acpCon)
	if err != nil {
		log.Println(err)
		switch strings(err.Error()) {
		case ""
		}
		//construct fail message as well
		//depends on err reply message (rule set)

		return
	}
	//return base on
	//should response based on cmd type??

	return
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
