package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	socketPath := "/tmp/unixgram_server.sock"
	os.Remove(socketPath)

	addr, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 創建 unixgram 伺服器
	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		fmt.Println("Error creating server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Unixgram server listening on", socketPath)

	buf := make([]byte, 1024)
	for {
		// 接收數據
		bufn, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			continue
		}
		fmt.Printf("Received: %s from %s\n", string(buf[:bufn]), addr.String())

		// 回應數據
		_, err = conn.WriteTo([]byte("Hello from Unixgram server"), addr)
		if err != nil {
			fmt.Println("Error writing to connection:", err)
		}
	}
}
