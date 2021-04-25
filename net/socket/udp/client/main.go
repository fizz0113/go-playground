package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("Hello UDP server"))
	buf := make([]byte, 1024)
	n, _, _ := conn.ReadFromUDP(buf)
	fmt.Println("Received:", string(buf[:n]))
}
