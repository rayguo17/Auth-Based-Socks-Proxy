package main

import (
	"context"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"github.com/rayguo17/go-socks/cmd/protocol/socks"
	"github.com/rayguo17/go-socks/util"
	"github.com/rayguo17/go-socks/util/logger"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Too less Arguments")
		return
	}
	socksListener, err := net.Listen("tcp", "0.0.0.0:"+os.Args[1])
	if err != nil {
		return
	}
	log.Println("socks running")
	go handleSocksListener(socksListener)
	transportListener, err := net.Listen("tcp", "0.0.0.0:"+os.Args[2])
	if err != nil {
		return
	}
	log.Println("transport running")
	go handleTransportListner(transportListener)
	select {}
}
func handleSocksListener(listener net.Listener) {
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Accept socks connection")
		go handleSocks(conn)
	}
}
func handleSocks(con net.Conn) {
	socksBuf := make([]byte, 512)

	_, err := con.Read(socksBuf)
	if err != nil {
		log.Println(err)
		return
	}
	con.Write([]byte{5, 0})
	cmdBuf := make([]byte, 512)
	cmdLen, err := con.Read(cmdBuf)

	lightConn, err := net.Dial("tcp", "127.0.0.1:"+os.Args[3])
	//add authentication
	lightConn.Write(cmdBuf[:cmdLen])
	lightResp := make([]byte, 512)

	lightConn.Read(lightResp)
	if lightResp[0] != 1 {
		//fail
		log.Println("got negative feedback from transport server")
		return
	}
	log.Print("get request ")
	pp.Println(cmdBuf[:cmdLen])
	con.Write([]byte{
		5, 0, 0, 1, 127, 0, 0, 1, 13, 88}) //hard code 5000
	//success start copy
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		io.Copy(con, lightConn)
		cancel()
	}()
	go func() {
		io.Copy(lightConn, con)
		cancel()
	}()
	<-ctx.Done()
}
func handleTransportListner(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Accepted transport connection")
		go handleTransport(conn)
	}
}
func handleTransport(conn net.Conn) {
	buf := make([]byte, 512)
	bufLen, _ := conn.Read(buf)
	source, err := socks.FromByte(buf[:bufLen], socks.ClientCommand)
	if err != nil {
		conn.Write([]byte{
			0,
		})
		log.Println("error socks command recognizing")
		return
	}
	if _, ok := source.(*socks.ClientCmd); !ok {
		logger.Debug.Println("socks authReq mapping failed returning")
		conn.Write([]byte{
			0,
		})
		return
	}
	clientCmd := source.(*socks.ClientCmd)
	var DstAddr util.Address
	switch clientCmd.Atyp {
	case socks.DomainAddress:
		DstAddr = util.NewDomainAddr(clientCmd.DstAddr, clientCmd.DstPort)
	case socks.Ipv4Address:
		DstAddr = util.NewIpv4Addr(clientCmd.DstAddr, clientCmd.DstPort)
	}
	log.Println("Get request", DstAddr.String())
	targetCon, err := net.Dial("tcp", DstAddr.String())
	if err != nil {
		conn.Write([]byte{
			0,
		})
		log.Println(err)
		return

	}
	conn.Write([]byte{
		1,
	})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		io.Copy(conn, targetCon)
		cancel()
	}()
	go func() {
		io.Copy(targetCon, conn)
		cancel()
	}()
	<-ctx.Done()
	return
}
