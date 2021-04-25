package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("Hello TCP server"))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("Received:", string(buf[:n]))
}
