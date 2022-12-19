package light

import (
	"context"
	"log"
	"net"

	pt "git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/util/logger"
	"github.com/rayguo17/go-socks/util/protocol/light"
	"gitlab.com/yawning/obfs4.git/common/drbg"
	"gitlab.com/yawning/obfs4.git/transports/obfs4"
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
			go acceptHandler(conn, system.LightConfig)
		}
	}()
	return nil
}

func acceptHandler(conn net.Conn, lightConfig config.LightConfig) {
	//handle obfs4 connection first
	t := obfs4.Transport{}
	seed, err := drbg.NewSeed()
	if err != nil {
		logger.Debug.Fatal(err)
	}
	pArgs := &pt.Args{
		"node-id":     []string{lightConfig.NodeID},
		"private-key": []string{lightConfig.PrivateKey},
		"drbg-seed":   []string{seed.Hex()},
		"iat-mode":    []string{"0"},
	}
	f, err := t.ServerFactory("./", pArgs)
	if err != nil {
		logger.Debug.Fatal(err)
	}
	addrStr := conn.RemoteAddr().String()
	name := f.Transport().Name()
	remote, err := f.WrapConn(conn)

	//
	if err != nil {
		log.Printf("%s(%s) - handshake failed: %s", name, addrStr, err)
		return
	}
	log.Printf("%s - handshake success", addrStr)
	//start light authentication.
	authBuf := make([]byte, 512)
	authLen, err := remote.Read(authBuf)
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
	acpCon, err := Authentication(ar, remote)

	//pp.Println(acpCon)
	if err != nil {
		logger.Debug.Println(err)
		_, err = remote.Write([]byte{1})
		if err != nil {
			logger.Debug.Println(err)
			return
		}
		return
	}
	defer acpCon.ProtocolClose()
	_, err = remote.Write([]byte{0})
	if err != nil {
		logger.Debug.Println(err)
		return
	}
	cmdBuf := make([]byte, 100)
	cmdLen, err := remote.Read(cmdBuf)
	if err != nil {
		logger.Debug.Println(err)
		remote.Write([]byte{1})
		return
	}
	cmd, err := light.BuildCmd(cmdBuf[:cmdLen])
	err = ConnectHandle(cmd, acpCon)
	if err != nil {
		logger.Debug.Println(err)
		_, err = remote.Write([]byte{1})
		//construct fail message as well
		//depends on err reply message (rule set)
		return
	}
	_, err = remote.Write([]byte{0})
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
