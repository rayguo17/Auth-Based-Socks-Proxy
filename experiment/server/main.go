package main

import (
	"fmt"
	pt "git.torproject.org/pluggable-transports/goptlib.git"
	"gitlab.com/yawning/obfs4.git/common/drbg"
	"gitlab.com/yawning/obfs4.git/transports/obfs4"
	"io"
	"log"
	"net"
	"sync"
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

const PublicKey = "3d92815f18a1cdf9911dfb0ba49c2586927e61962e94fb19a3a630cc10ba0528"
const PrivateKey = "70FA07A923B0FDFFD797F59AC4F4DFE2B516D5F6ECCB33EE29BCFC94989F2E67"

func acceptHandler(conn net.Conn) {
	t := obfs4.Transport{}
	log.Println("connection accepted")
	seed, err := drbg.NewSeed()
	if err != nil {
		log.Fatal(err)
	}
	//seed.Hex()
	pArgs := &pt.Args{
		"node-id":     []string{"A868303126987902D51F2B6F06DD90038C45B119"},
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
	fmt.Println("Handshake success!")
	receiveConnection, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatal(err)
	}
	defer receiveConnection.Close()
	if err = copyLoop(receiveConnection, remote); err != nil {
		log.Printf("%s(%s) - closed connection: %s", name, addrStr, err)
	} else {
		log.Printf("%s(%s) - closed connection", name, addrStr)
	}
	//pp.Println(remote)

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