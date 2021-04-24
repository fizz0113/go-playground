package main

import (
	"fmt"
	"net"
	"net/rpc"
)

func main() {
	// Connect to the RPC server
	conn, err := net.Dial("tcp", "localhost:8080") // 替換為伺服器地址
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create a new RPC client
	client := rpc.NewClient(conn)

	// Call the method on the server
	req, res := "ian", ""
	err = client.Call("XXXService.Greet", req, &res)
	if err != nil {
		panic(err)
	}

	// Print the server response
	fmt.Printf("Server response: %s\n", res)
}
