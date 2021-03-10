package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	grpc_backoff "google.golang.org/grpc/backoff"
	grpc_insecure "google.golang.org/grpc/credentials/insecure"

	"grpc/stream/prototypes"
)

var (
	serverHostName = flag.String("serverHostName", "localhost", "Specify ip/hostname")
	serverHostPort = flag.Uint("serverHostPort", 6666, "Specify host port")
	topics         = flag.String("topics", "topic1,topic2", "Specify topics, use ',' to split")
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

	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true), // blocking until connect server
		),
		grpc.WithBlock(), // blocking until connect server
		// grpc.WithContextDialer(makeClientDialer), // blocking until TCP dial complete
		grpc.WithTransportCredentials(grpc_insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: grpc_backoff.Config{
				BaseDelay:  500 * time.Millisecond,
				Multiplier: 1.1,
				Jitter:     0.1,
				MaxDelay:   120 * time.Second,
			},
			MinConnectTimeout: time.Millisecond,
		}),
		// client middleware
		// grpc.WithStreamInterceptor(
		// ),
		// grpc.WithChainUnaryInterceptor(
		// ),
	}
	host := fmt.Sprintf("%s:%d", *serverHostName, *serverHostPort)
	clientConn, err := grpc.Dial(host, opts...)
	if err != nil {
		panic(err)
	}
	defer clientConn.Close()

	brokerServiceClient := prototypes.NewBrokerServiceClient(clientConn)

	cfg := &prototypes.SubscriberConfig{
		Topics:      strings.Split(*topics, ","),
		ChannelSize: 128,
	}
	subscriber, err := brokerServiceClient.Subscribe(ctx, cfg)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-subscriber.Context().Done():
			return
		case <-time.After(time.Second):
			msg, err := subscriber.Recv()
			if err != nil {
				panic(err)
			}
			log.Printf("Received a message: %+v", msg)
		}
	}
}
