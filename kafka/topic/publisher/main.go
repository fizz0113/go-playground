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

func publishMessage(conn *kafka.Conn) {
	msg := kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339Nano)),
		Value: []byte("Hello Kafka " + time.Now().String()),
	}
	_, err := conn.WriteMessages(msg)
	if err != nil {
		log.Println("Failed to write a message:", err)
		return
	}
	log.Printf("Sent a message %q", msg)
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

	// 獲取控制器
	controller, err := conn.Controller()
	panicWhenErrorNotNil("Failed to get controller:", err)
	log.Println("get controller succeed")

	// 透過獲取的的控制資訊連接 kafka
	controllerHost := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	controllerConn, err := kafka.Dial("tcp", controllerHost)
	panicWhenErrorNotNil("Failed to connect controller:", err)
	defer controllerConn.Close()
	log.Println(controllerHost, "controller connected")

	// 建立 topic
	topicCfg := kafka.TopicConfig{
		Topic:              *topic,
		NumPartitions:      2,
		ReplicationFactor:  1,
		ReplicaAssignments: []kafka.ReplicaAssignment{},
		ConfigEntries:      []kafka.ConfigEntry{},
	}
	err = controllerConn.CreateTopics(topicCfg)
	panicWhenErrorNotNil("Failed to create topic:", err)
	defer controllerConn.DeleteTopics(*topic)
	log.Println(*topic, "created")

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
		case <-time.After(time.Second):
			publishMessage(leader)
		}
	}
}
