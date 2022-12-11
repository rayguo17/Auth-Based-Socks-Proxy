package main

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("receiver ready")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	buf := make([]byte, 512)
	bufLen, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(buf[:bufLen])
	pp.Println(string(buf[:bufLen]))
}
