package light

import (
	"context"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/util/logger"
	"github.com/rayguo17/go-socks/util/protocol/light"
	"net"
)

func ListenStart(system *config.System, cancelFunc context.CancelFunc) error {
	listener, err := net.Listen("tcp", "0.0.0.0:"+system.GetLightPort())
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
	authBuf := make([]byte, 512)
	authLen, err := conn.Read(authBuf)
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	ar, err := light.BuildAR(authBuf[:authLen])
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	//pp.Println(ar)
	acpCon, err := Authentication(ar, conn)

	//pp.Println(acpCon)
	if err != nil {
		logger.Debug.Println(err)
		_, err = conn.Write([]byte{1})
		if err != nil {
			logger.Debug.Println(err)
			return
		}
		return
	}
	defer acpCon.ProtocolClose()
	_, err = conn.Write([]byte{0})
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	cmdBuf := make([]byte, 100)
	cmdLen, err := conn.Read(cmdBuf)
	if err != nil {
		logger.Debug.Println(err)
		conn.Write([]byte{1})
		return
	}
	cmd, err := light.BuildCmd(cmdBuf[:cmdLen])
	err = ConnectHandle(cmd, acpCon)
	if err != nil {
		logger.Debug.Println(err)
		_, err = conn.Write([]byte{1})
		//construct fail message as well
		//depends on err reply message (rule set)
		return
	}
	_, err = conn.Write([]byte{0})
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	err = acpCon.ExecuteBegin()
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	return
}
