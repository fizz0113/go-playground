package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer l.Close()

	fmt.Println("TCP server listening on port 8080...")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			fmt.Println("Received:", string(buf[:n]))
			c.Write([]byte("Hello from TCP server"))
		}(conn)
	}
}
