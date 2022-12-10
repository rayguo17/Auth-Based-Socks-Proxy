package main

import (
	"fmt"
	pt "git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/k0kubun/pp/v3"
	"gitlab.com/yawning/obfs4.git/common/drbg"
	"gitlab.com/yawning/obfs4.git/transports/obfs4"
	"log"
	"net"
)

func main() {
	port := "10009"
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go acceptHandler(conn)
	}

}

const PublicKey = "c2ab50c38a7d19103066c3b7612e7ac286c2aec3985ec2bcc3b7651551c282c2900cc395c3aec298c3a47c30c2a867"
const PrivateKey = "3a56dcd29a6bc5fe6a2f36534e612c701d0c0e8196dd6a700bb180f6d3dc8bdb"

func acceptHandler(conn net.Conn) {
	t := obfs4.Transport{}
	log.Println("connection accepted")
	seed, err := drbg.NewSeed()
	if err != nil {
		log.Fatal(err)
	}
	//seed.Hex()
	pArgs := &pt.Args{
		"node-id":     []string{"0077BCBA7244DB3E6A5ED2746E86170066684887"},
		"private-key": []string{PrivateKey},
		"drbg-seed":   []string{seed.Hex()},
		"iat-mode":    []string{"0"},
	}
	f, err := t.ServerFactory("./", pArgs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("wraping connection")
	addrStr := conn.RemoteAddr().String()
	name := f.Transport().Name()
	remote, err := f.WrapConn(conn)

	if err != nil {
		log.Printf("%s(%s) - handshake failed: %s", name, addrStr, err)
		return
	}
	pp.Println(remote)
	fmt.Println("Connection success!")
}
