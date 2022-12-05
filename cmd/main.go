package main

import (
	"fmt"
	"github.com/rayguo17/go-socks/Backdoor"
	"github.com/rayguo17/go-socks/socks"
	"github.com/rayguo17/go-socks/user"
	"log"
	"net"
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
		_, err = conn.Write([]byte{1, 0})
		if err != nil {
			log.Println(err)
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
	defer acpCon.ProtocolClose() //close on return

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
		log.Println("socks authReq mapping failed returning")
		return
	}
	clientCmd := source.(*socks.ClientCmd)
	//pp.Println(clientCmd)
	//commandHandle -> handleConnect -> con.ConnectCmd
	//fmt.Println("going to handle command")
	err = socks.CommandHandle(clientCmd, acpCon)

	if err != nil {
		log.Println(err)
		_, err = conn.Write([]byte{5, 1, 0, 1, 1, 2, 3, 4, 1, 2})
		//construct fail message as well
		//depends on err reply message (rule set)
		return
	}
	//fmt.Println("Command handle success")
	cmdResp, err := socks.ConstructResp(acpCon, socks.ServerResponse)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = conn.Write(cmdResp)
	if err != nil {
		log.Println(err)
		return
	}
	//return base on
	//should response based on cmd type??
	//Final: start Executor Go Routine, Start Data Transfer
	err = acpCon.ExecuteBegin()
	log.Println("Executing")
	if err != nil {
		log.Println(err)
		return
	}
	select {}
	return

}

//browser tends to send multiple tcp connections, so there will be parallel thread using same instance.
