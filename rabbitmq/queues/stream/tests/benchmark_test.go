package tests

import (
	"fmt"
	"testing"

	streamAmqp "github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

var (
	globalConsumerRecvChannel chan *streamAmqp.Message
	globalConsumerEnv         *stream.Environment
	globalConsumer            *stream.Consumer
)

func initConsumer() {
	globalConsumerRecvChannel = make(chan *streamAmqp.Message, 102400)
	// 建立連線環境
	envOpts := stream.NewEnvironmentOptions().
		SetHost("rabbitmq").
		SetPort(5552).
		SetUser("admin").
		SetPassword("admin").
		SetVHost("vhost_admin")
	env, err := stream.NewEnvironment(envOpts)
	if err != nil {
		panic("failed to create environment:" + err.Error())
	}
	globalConsumerEnv = env

	// 建立stream
	streamName := "test-stream"
	streamOpts := stream.NewStreamOptions().
		SetMaxSegmentSizeBytes(stream.ByteCapacity{}.MB(1)).
		SetMaxLengthBytes(stream.ByteCapacity{}.MB(2))
	env.DeleteStream(streamName)
	err = env.DeclareStream(streamName, streamOpts)
	if err != nil {
		panic("Failed to declare stream:" + err.Error())
	}

	handleReceivedMessage := func(consumerContext stream.ConsumerContext, msg *streamAmqp.Message) {
		globalConsumerRecvChannel <- msg
	}
	// 建立consumer
	consumerOpts := stream.NewConsumerOptions().
		// SetClientProvidedName("my_consumer").           // connection name
		// SetConsumerName("my_consumer").                 // set a consumer name
		SetOffset(stream.OffsetSpecification{}.First()) // start consuming from the beginning
	consumer, err := env.NewConsumer(streamName, handleReceivedMessage, consumerOpts)
	if err != nil {
		panic("Failed to create consumer:" + err.Error())
	}
	globalConsumer = consumer
}

var (
	globalProducerEnv *stream.Environment
	globalProducer    *stream.Producer
)

func initProducer() {
	// 建立連線環境
	envOpts := stream.NewEnvironmentOptions().
		SetHost("rabbitmq").
		SetPort(5552).
		SetUser("admin").
		SetPassword("admin").
		SetVHost("vhost_admin")
	env, err := stream.NewEnvironment(envOpts)
	if err != nil {
		panic("failed to create environment:" + err.Error())
	}
	globalProducerEnv = env

	// 建立stream
	streamName := "test-stream"

	// 建立producer
	producerOpts := stream.NewProducerOptions().
		SetSubEntrySize(100).
		SetCompression(stream.Compression{}.Gzip()).
		SetProducerName("producer")
	producer, err := env.NewProducer(streamName, producerOpts)
	if err != nil {
		panic("Failed to create producer: " + err.Error())
	}
	globalProducer = producer
}

func init() {
	initConsumer()
	initProducer()
}

func Benchmark_StreamQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msg := fmt.Sprintf("msg-%d", i)
		err := globalProducer.Send(streamAmqp.NewMessage([]byte(msg)))
		if err != nil {
			b.Fatal("Failed to send message: ", err)
		}
		<-globalConsumerRecvChannel
	}
	b.Logf("Processed %d messages", b.N)
}

/*
❯ go test -bench=. -v  .
Alias tip: got -bench=. -v  .
goos: linux
goarch: amd64
pkg: rabbitmq/stream/tests
cpu: Intel(R) Core(TM) i5-10210U CPU @ 1.60GHz
Benchmark_StreamQueue
    benchmark_test.go:102: Processed 1 messages
    benchmark_test.go:102: Processed 100 messages
    benchmark_test.go:102: Processed 550 messages
Benchmark_StreamQueue-4              550           2123572 ns/op
PASS
ok      rabbitmq/stream/tests   1.599s
*/
