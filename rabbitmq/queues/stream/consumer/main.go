package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	streamAmqp "github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
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

func handleReceivedMessage(consumerContext stream.ConsumerContext, msg *streamAmqp.Message) {
	fmt.Printf("data: %v \n", string(msg.GetData()))
	consumerContext.Consumer.StoreOffset() // 👈 確保 offset 被存起來，不會掉包
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

	// 建立consumer
	consumerOpts := stream.NewConsumerOptions().
		// SetClientProvidedName("my_consumer").           // connection name
		// SetConsumerName("my_consumer").                 // set a consumer name
		SetOffset(stream.OffsetSpecification{}.First()) // start consuming from the beginning
	consumer, err := env.NewConsumer(*queueName, handleReceivedMessage, consumerOpts)
	panicWhenErrorNotNil("Failed to create consumer: %v", err)
	defer consumer.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-consumer.NotifyClose():
			return
		}
	}
}
