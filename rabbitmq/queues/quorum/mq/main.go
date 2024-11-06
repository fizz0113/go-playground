package main

import (
	"context"
	"flag"
	"fmt"
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
	queueName        = flag.String("queueName", "test_quorum_queue", "Specify queue name")
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

	q, err := ch.QueueDeclare(
		*queueName, // name
		true,       // durable must be true
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		amqp.Table{ // arguments
			"x-queue-type": "quorum",
			// "x-expires":                   10000, // Number
			// "x-message-ttl":               10, // Number
			// "x-overflow":                  "", // String
			// "x-single-active-consumer":    false, // Boolean
			// "x-dead-letter-exchange":      "", // String
			// "x-dead-letter-routing-key":   "", // String
			// "x-max-length":                10, // Number
			// "x-max-length-bytes":          10, // Number
			// "x-delivery-limit":            10, // Number
			// "x-quorum-initial-group-size": 10, // Number
			// "x-dead-letter-strategy":      "", // String
			// "x-queue-leader-locator":      "", // String
		},
	)
	panicWhenErrorNotNil("Failed to declare quorum queue:", err)
	defer ch.QueueDelete(q.Name, false, false, false)

	<-ctx.Done()
}
