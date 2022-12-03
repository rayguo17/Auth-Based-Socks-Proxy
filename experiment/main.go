package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "www.baidu.com:80")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("dial success")
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	//reader := bufio.NewReader(conn)
	var httpBody string
	var httpSize int
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		httpSize += len
		httpBody += string(buf[:len])
		if len < 512 {
			break
		}

	}
	log.Println(httpBody, httpSize)
}
