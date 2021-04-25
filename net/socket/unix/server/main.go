package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	socketPath := "/tmp/unix.sock"
	os.Remove(socketPath)

	addr, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	l, err := net.ListenUnix("unix", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer l.Close()

	fmt.Println("Unix socket server listening...")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			bufn, err := c.Read(buf)
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			fmt.Println("Received:", string(buf[:bufn]))
			c.Write([]byte("Hello from Unix socket server"))
		}(conn)
	}
}
