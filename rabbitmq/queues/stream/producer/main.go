package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx
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
	panicWhenErrorNotNil("Failed to create environment: ", err)
	defer env.Close()

	// 建立producer
	producerOpts := stream.NewProducerOptions().
		SetSubEntrySize(100).
		SetCompression(stream.Compression{}.Gzip()).
		SetProducerName("producer")
	producer, err := env.NewProducer(*queueName, producerOpts)
	panicWhenErrorNotNil("Failed to create producer: ", err)
	defer producer.Close()

	chPublishConfirm := producer.NotifyPublishConfirmation()

	// 發送 10 筆訊息
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("msg-%d", i)
		err = producer.Send(streamAmqp.NewMessage([]byte(msg)))
		if err != nil {
			panicWhenErrorNotNil("Failed to send message: ", err)
		}
		log.Println("Sent: ", msg)
	}

	for {
		select {
		case confirmed := <-chPublishConfirm:
			for _, msg := range confirmed {
				if msg.IsConfirmed() {
					fmt.Printf("message %s stored \n", msg.GetMessage().GetData())
				} else {
					fmt.Printf("message %s failed \n", msg.GetMessage().GetData())
				}
			}
		case <-time.After(time.Second):
			return
		}
	}
}
