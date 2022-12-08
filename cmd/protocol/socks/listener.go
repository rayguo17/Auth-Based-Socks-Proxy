package socks

import (
	"context"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/manager/connection/socks"
	"github.com/rayguo17/go-socks/util/logger"
	"net"
)

func ListenStart(system *config.System, cancelFunc context.CancelFunc) error {

	listener, err := net.Listen("tcp", "0.0.0.0"+":"+system.GetSocksPort())
	if err != nil {

		return err
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Debug.Println(err)
				cancelFunc()
				return
			}
			go acceptHandler(conn)
		}
	}()
	return nil
}

func acceptHandler(conn net.Conn) {
	buf := make([]byte, 512)
	initLen, err := conn.Read(buf)
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	//fmt.Printf("%d bytes received!\n", initLen)
	source, err := FromByte(buf[:initLen], HandShakeRequest)
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	if _, ok := source.(*HandshakeReq); !ok {
		logger.Debug.Println("socks handshake mapping failed returning")
		return
	}
	handShakeReq := source.(*HandshakeReq)
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
			logger.Debug.Println(err)
			return
		}
	case 2:
		_, err = conn.Write([]byte{5, 2})
		if err != nil {
			logger.Debug.Println(err)
			return
		}
	}
	//3. authentication phase (skip for now)
	var acpCon *socks.AcpCon
	if chosenAuthMethod == 2 {
		authBuf := make([]byte, 512)
		authLen, err := conn.Read(authBuf)
		if err != nil {
			logger.Debug.Println(err)
			return
		}
		//register tmp to received feedback.

		source, err = FromByte(authBuf[:authLen], AuthRequest)
		if err != nil {
			logger.Debug.Println(err)
			return
		}
		if _, ok := source.(*AuthReq); !ok {
			logger.Debug.Println("socks authReq mapping failed returning")
			return
		}
		authReq := source.(*AuthReq)
		acpCon, err = Authenticate(authReq, conn)
		if err != nil {
			logger.Debug.Println(err)
			conn.Write([]byte{1, 1})
			return
		}
		_, err = conn.Write([]byte{1, 0})
		if err != nil {
			logger.Debug.Println(err)
			return
		}
	} else {
		//create anonymous account route
		acpCon, err = Anonymous(conn)
		if err != nil {
			logger.Debug.Println(err)
			conn.Write([]byte{1, 1})
			return
		}
	}
	defer acpCon.ProtocolClose() //close on return

	//TODO://4. request phase
	cmdBuf := make([]byte, 512)
	cmdLen, err := conn.Read(cmdBuf)
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	source, err = FromByte(cmdBuf[:cmdLen], ClientCommand)
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	if _, ok := source.(*ClientCmd); !ok {
		logger.Debug.Println("socks authReq mapping failed returning")
		return
	}
	clientCmd := source.(*ClientCmd)
	//pp.Println(clientCmd)
	//commandHandle -> handleConnect -> con.ConnectCmd
	//fmt.Println("going to handle command")
	err = CommandHandle(clientCmd, acpCon)

	if err != nil {
		logger.Debug.Println(err)
		_, err = conn.Write([]byte{5, 1, 0, 1, 1, 2, 3, 4, 1, 2})
		//construct fail message as well
		//depends on err reply message (rule set)
		return
	}
	//fmt.Println("Command handle success")
	cmdResp, err := ConstructResp(acpCon, ServerResponse)
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	_, err = conn.Write(cmdResp)
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	//return base on
	//should response based on cmd type??
	//Final: start Executor Go Routine, Start Data Transfer
	err = acpCon.ExecuteBegin()
	//log.Println("Executing")
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	//
	return

}
