package main

import (
	"fmt"
	"net"
)

func main() {
	socketPath := "/tmp/unix.sock"

	serverAddr, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.DialUnix("unix", nil, serverAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("Hello Unix socket server"))
	buf := make([]byte, 1024)
	bufn, _ := conn.Read(buf)
	fmt.Println("Received:", string(buf[:bufn]))
}
