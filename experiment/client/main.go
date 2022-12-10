package main

import (
	"fmt"
	pt "git.torproject.org/pluggable-transports/goptlib.git"
	"gitlab.com/yawning/obfs4.git/transports/obfs4"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"sync"
)

func main() {
	socksPort := "5000"
	socksListener, err := net.Listen("tcp", "0.0.0.0:"+socksPort)
	if err != nil {
		return
	}
	log.Println("socks running")
	go handleSocksListener(socksListener)
	fmt.Println("hello")

	//sessionKey, err := ntor.NewKeypair(true)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//pp.Println(sessionKey)
	//hexString := "d97a9dcecc6ab9f22b1d8b081c51db9befe0e238d1209d87c41d62c269823643"
	//publicKey, err := ntor.PublicKeyFromHex(hexString)
	//pp.Println(publicKey)
	////fmt.Println(dir)
	//nodeId, err := ntor.NodeIDFromHex("A868303126987902D51F2B6F06DD90038C45B119")
	//pp.Println(nodeId)
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

func handleSocks(conn net.Conn) {
	t := obfs4.Transport{}
	//dir, err := pt.MakeStateDir()
	f, err := t.ClientFactory("./")
	if err != nil {
		log.Println(err)
		return
	}
	socksBuf := make([]byte, 512)

	_, err = conn.Read(socksBuf)
	if err != nil {
		log.Println(err)
		return
	}
	conn.Write([]byte{5, 0})
	cmdBuf := make([]byte, 512)
	cmdLen, err := conn.Read(cmdBuf)
	//pp.Println(cmdBuf[:cmdLen])
	//clientArgs
	sockArgs := &pt.Args{
		"node-id": []string{
			"A868303126987902D51F2B6F06DD90038C45B119",
		},
		"public-key": []string{
			"d97a9dcecc6ab9f22b1d8b081c51db9befe0e238d1209d87c41d62c269823643",
		},
		"iat-mode": []string{
			"0",
		},
	}

	name := f.Transport().Name()
	//fmt.Println(name)
	//addrStr := "www.baidu.com:443"
	dialFn := proxy.Direct.Dial
	args, err := f.ParseArgs(sockArgs)
	if err != nil {
		log.Println(err)
	}
	//pp.Println(args)
	addrStr := "127.0.0.1:10009"
	remote, err := f.Dial("tcp", addrStr, dialFn, args)
	if err != nil {
		log.Printf("%s(%s) - outgoing connection failed: %s", name, addrStr, err)
		//_ = socksReq.Reply(socks5.ErrorToReplyCode(err))
		return
	}
	defer remote.Close()

	//pp.Println(cmdBuf[:cmdLen])
	conn.Write([]byte{
		5, 0, 0, 1, 127, 0, 0, 1, 13, 88})
	if err = copyLoop(conn, remote); err != nil {
		log.Printf("%s(%s) - closed connection: %s", name, addrStr, err)
	}

}

func copyLoop(a net.Conn, b net.Conn) error {
	// Note: b is always the pt connection.  a is the SOCKS/ORPort connection.
	errChan := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer b.Close()
		defer a.Close()
		_, err := io.Copy(b, a)
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		defer a.Close()
		defer b.Close()
		_, err := io.Copy(a, b)
		errChan <- err
	}()

	// Wait for both upstream and downstream to close.  Since one side
	// terminating closes the other, the second error in the channel will be
	// something like EINVAL (though io.Copy() will swallow EOF), so only the
	// first error is returned.
	wg.Wait()
	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}
