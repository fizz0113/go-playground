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
	exchangeName     = flag.String("exchangeName", "test_headers_exchange", "Specify exchange name")
	match            = flag.String("match", "all", "Specify 'all' or 'any'")
)

func panicWhenErrorNotNil(msg string, err error) {
	if err != nil {
		panic(msg + err.Error())
	}
}

func handleReceivedMessage(msg amqp.Delivery) {
	log.Printf("Received a message: %q", msg.Body)
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
		"",         // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		amqp.Table{ // arguments
			// "x-expires":                 10000, // Number: ms
			// "x-message-ttl":             10,    // Number:
			// "x-overflow":                "",    // String
			// "x-single-active-consumer":  false, // Boolean
			// "x-dead-letter-exchange":    "",    // String
			// "x-dead-letter-routing-key": "",    // String
			// "x-max-length":              10,    // Number
			// "x-max-length-bytes":        10,    // Number
			// "x-max-priority":            10,    // Number
			// "x-queue-version":           10,    // Number
			// "x-queue-master-locator":    "",    // String
		},
	)
	panicWhenErrorNotNil("Failed to declare queue:", err)
	defer ch.QueueDelete(q.Name, false, false, false)
	log.Println("declare queue", q.Name)

	err = ch.QueueBind(
		q.Name,        // queue name
		"",            // default route key
		*exchangeName, // bind exchange
		false,         // no wait
		amqp.Table{ // args
			"x-match": *match,
			"format":  "pdf",
			"type":    "report",
		},
	)
	panicWhenErrorNotNil("Failed to declare queue:", err)
	defer ch.QueueUnbind(q.Name, "", *exchangeName, nil)
	log.Println("bind queue", q.Name, "to exchange", *exchangeName)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	panicWhenErrorNotNil("Failed to register a consumer:", err)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			handleReceivedMessage(msg)
		}
	}
}
