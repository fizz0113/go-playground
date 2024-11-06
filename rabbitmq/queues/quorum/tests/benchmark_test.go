package tests

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RABBITMQ_URL = "amqp://admin:admin@rabbitmq:5672/vhost_admin"
	QUEUE_NAME   = "test_quorum_queue_benchmark"
)

func Benchmark_QuormQueue(b *testing.B) {
	conn, err := amqp.Dial(RABBITMQ_URL)
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		b.Fatal(err)
	}
	defer ch.Close()

	args := amqp.Table{
		"x-queue-type": "quorum",
	}

	q, err := ch.QueueDeclare(QUEUE_NAME, true, false, false, false, args)
	if err != nil {
		b.Fatal(err)
	}
	defer ch.QueueDelete(q.Name, false, false, false)

	ch.QueuePurge(q.Name, false)

	msgs, err := ch.Consume(QUEUE_NAME, "", true, false, false, false, nil)
	if err != nil {
		b.Fatal(err)
	}

	func() {
		conn, _ := amqp.Dial(RABBITMQ_URL)
		defer conn.Close()
		ch, _ := conn.Channel()
		defer ch.Close()
		for i := 0; i < b.N; i++ {
			_ = ch.Publish("", QUEUE_NAME, false, false, amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte("benchmark pub-consume"),
				DeliveryMode: amqp.Persistent,
			})

			<-msgs
		}
	}()
	b.Logf("Processed %d messages", b.N)
}

/*
❯ go test -bench=. -v  .
Alias tip: got -bench=. -v  .
goos: linux
goarch: amd64
pkg: rabbitmq/quorum/tests
cpu: Intel(R) Core(TM) i5-10210U CPU @ 1.60GHz
Benchmark_QuormQueue
    benchmark_test.go:59: Processed 1 messages
    benchmark_test.go:59: Processed 18 messages
    benchmark_test.go:59: Processed 164 messages
    benchmark_test.go:59: Processed 214 messages
Benchmark_QuormQueue-4               214           6054331 ns/op
PASS
ok      rabbitmq/quorum/tests   2.411s
*/
