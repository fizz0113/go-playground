package main

import (
	"fmt"
	"net"
)

func main() {
	serverSocketPath := "/tmp/unixgram_server.sock"

	// 創建 unixgram 客戶端
	conn, err := net.Dial("unixgram", serverSocketPath)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Unixgram client connected to", serverSocketPath)

	// 發送數據
	message := "Hello Unixgram server"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to server:", err)
		return
	}
	fmt.Println("Sent:", message)

	// 接收伺服器回應
	buf := make([]byte, 1024)
	bufn, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}
	fmt.Println("Received:", string(buf[:bufn]))
}
