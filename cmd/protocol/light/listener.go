package light

import (
	"context"
	"github.com/k0kubun/pp/v3"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/util/logger"
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
	ar, err := BuildAR(authBuf[:authLen])
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	acpCon, err := Authentication(ar)
	pp.Println(acpCon)
}
