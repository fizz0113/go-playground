package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaHostName = flag.String("kafkaHostName", "kafka-node1", "Specify host name")
	kafkaHostPort = flag.Uint("kafkaHostPort", 9092, "Specify host port")
	topic         = flag.String("topic", "test_topic", "Specify topic")
)

func panicWhenErrorNotNil(msg string, err error) {
	if err != nil {
		panic(msg + err.Error())
	}
}

func receiveMessage(conn *kafka.Conn) {
	batch := conn.ReadBatch(1e3, 1e6)
	defer batch.Close()
	for {
		msg, err := batch.ReadMessage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Println("Failed to read a message:", err)
			return
		}
		log.Printf("Recv a message %q", msg)
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

	// 連線到 Kafka broker
	host := fmt.Sprintf("%s:%d", *kafkaHostName, *kafkaHostPort)
	conn, err := kafka.Dial("tcp", host)
	panicWhenErrorNotNil("Failed to connect kafka:", err)
	defer conn.Close()
	log.Println(host, "connected")

	// 檢查 topic partition 是否存在
	partitions, err := conn.ReadPartitions(*topic)
	panicWhenErrorNotNil("Failed to get partitions:", err)
	for _, p := range partitions {
		log.Printf("Topic=%s Partition=%d Leader=%v Replicas=%v ISR=%v", p.Topic, p.ID, p.Leader, p.Replicas, p.Isr)
	}

	// 連到 partition leader
	leaderHost := fmt.Sprintf("%s:%d", partitions[0].Leader.Host, partitions[0].Leader.Port)
	partition := 0
	leader, err := kafka.DialLeader(ctx, "tcp", leaderHost, *topic, partition)
	panicWhenErrorNotNil("Failed to connect partition leader:", err)
	defer leader.Close()
	log.Println(leaderHost, "partition leader connected")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			receiveMessage(leader)
		}
	}
}
