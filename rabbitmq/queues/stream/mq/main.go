package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	stream "github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

var (
	rabbitmqHostName = flag.String("rabbitmqHostName", "rabbitmq", "Specify rabbitmq host name")
	rabbitmqHostPort = flag.Int("rabbitmqHostPort", 5552, "Specify rabbitmq host port")
	rabbitmqUsername = flag.String("rabbitmqUsername", "admin", "Specify rabbitmq login username")
	rabbitmqPassword = flag.String("rabbitmqPassword", "admin", "Specify rabbitmq login password")
	rabbitmqVHost    = flag.String("rabbitmqVHost", "vhost_admin", "Specify rabbitmq vhost")
	queueName        = flag.String("queueName", "test_stream_queue", "Specify queue name")
)

func panicWhenErrorNotNil(msg string, err error) {
	if err != nil {
		panic(msg + err.Error())
	}
}

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

	// 建立連線環境
	envOpts := stream.NewEnvironmentOptions().
		SetHost(*rabbitmqHostName).
		SetPort(*rabbitmqHostPort).
		SetUser(*rabbitmqUsername).
		SetPassword(*rabbitmqPassword).
		SetVHost(*rabbitmqVHost)
	env, err := stream.NewEnvironment(envOpts)
	panicWhenErrorNotNil("failed to create environment:", err)
	defer env.Close()

	// 建立stream
	streamOpts := stream.NewStreamOptions().
		SetMaxSegmentSizeBytes(stream.ByteCapacity{}.MB(1)).
		SetMaxLengthBytes(stream.ByteCapacity{}.MB(2))
	env.DeleteStream(*queueName)
	err = env.DeclareStream(*queueName, streamOpts)
	panicWhenErrorNotNil("Failed to declare stream: ", err)
	defer env.DeleteStream(*queueName)

	<-ctx.Done()
}
