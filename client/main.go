package client

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	conn.Write([]byte{5, 1, 0})
	reader := bufio.NewReader(conn)
	buf1 := make([]byte, 512)
	len1, err := reader.Read(buf1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%d bytes received\n", len1)
	fmt.Println(buf1[:len1])

	conn.Write([]byte{5, 1, 0, 3, 119, 119, 119, 46, 98, 97, 105, 100, 117, 46, 99, 111, 109, 0, 80})

}
