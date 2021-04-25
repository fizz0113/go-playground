package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server listening on port 8080...")
	buf := make([]byte, 1024)
	for {
		bufn, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		fmt.Println("Received:", string(buf[:bufn]), "from", clientAddr)
		conn.WriteToUDP([]byte("Hello from UDP server"), clientAddr)
	}
}
