package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	// grpc_middleware_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	// grpc_middleware_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	// grpc_middleware_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	// grpc_middleware_prometheus "github.com/grpc-ecosystem/go-grpc-middleware/prometheus"

	"google.golang.org/grpc"
	// grpc_reflection "google.golang.org/grpc/reflection"

	"grpc/basic/prototypes"
	"grpc/basic/server/hello_service"
)

var (
	serverHostName = flag.String("serverHostName", "localhost", "Specify ip/hostname")
	serverHostPort = flag.Uint("serverHostPort", 6666, "Specify host port")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	defer close(c)
	go func() {
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		// grpc_middleware_recovery.StreamServerInterceptor(),
		// grpc_middleware_ctxtags.StreamServerInterceptor(),
		// grpc_middleware_auth.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		// grpc_middleware_recovery.UnaryServerInterceptor(),
		// grpc_middleware_ctxtags.UnaryServerInterceptor(),
		// grpc_middleware_auth.UnaryServerInterceptor(authenticator),
		// grpc_opentracing.UnaryServerInterceptor(),
		// grpc_middleware_prometheus.UnaryServerInterceptor,
		// grpc_zap.UnaryServerInterceptor(zapLogger),
		)),
	}
	server := grpc.NewServer(opts...)
	prototypes.RegisterHelloServiceServer(server, hello_service.HelloService{})
	// grpc_reflection.Register(server)

	host := fmt.Sprintf("%s:%d", *serverHostName, *serverHostPort)
	l, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	log.Printf("Start server on %s", host)
	go server.Serve(l)
	<-ctx.Done()
}
