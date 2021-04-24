package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type XXXService struct{}

func (i *XXXService) Greet(req string, res *string) error {
	*res = "Hello, " + req
	return nil
}

func main() {
	// Define the TCP address to listen on
	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	// Create a TCP listener
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	// Create a new RPC server
	rpcServer := rpc.NewServer()

	// Register the service
	rpcServer.RegisterName("XXXService", new(XXXService))

	// Accept incoming connections and serve requests
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print("rpc.Serve: accept:", err.Error())
			continue
		}
		codec := jsonrpc.NewServerCodec(conn)
		rpcServer.ServeCodec(codec)
	}
}
