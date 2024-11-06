package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitmqHost     = flag.String("rabbitmqHost", "rabbitmq:5672", "Specify rabbitmq host")
	rabbitmqUsername = flag.String("rabbitmqUsername", "admin", "Specify rabbitmq login username")
	rabbitmqPassword = flag.String("rabbitmqPassword", "admin", "Specify rabbitmq login password")
	rabbitmqVHost    = flag.String("rabbitmqVHost", "vhost_admin", "Specify rabbitmq vhost")
	exchangeName     = flag.String("exchangeName", "test_fanout_exchange", "Specify exchange name")
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

	rabbitmqUrl := fmt.Sprintf("amqp://%s:%s@%s/%s",
		*rabbitmqUsername,
		*rabbitmqPassword,
		*rabbitmqHost,
		*rabbitmqVHost,
	)

	conn, err := amqp.Dial(rabbitmqUrl)
	panicWhenErrorNotNil("Failed to connect to rabbitmq:", err)
	defer conn.Close()

	ch, err := conn.Channel()
	panicWhenErrorNotNil("Failed to open a channel:", err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		*exchangeName, // name
		"fanout",      // kind
		false,         // durable
		false,         // auto delete
		false,         // internal
		false,         // no waite
		amqp.Table{    // args
			// "alternate-exchange": "",
		},
	)
	panicWhenErrorNotNil("Failed to declare a exchange:", err)
	defer ch.ExchangeDelete(*exchangeName, false, false)
	log.Println("declare exchange", *exchangeName)

	<-ctx.Done()
}
