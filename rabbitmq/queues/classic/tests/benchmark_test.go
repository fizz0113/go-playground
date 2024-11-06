package tests

import (
	"sync"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RABBITMQ_URL = "amqp://admin:admin@rabbitmq:5672/vhost_admin"
	QUEUE_NAME   = "test_classic_queue_benchmark"
)

func Benchmark_ClassicQueue(b *testing.B) {
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

	q, err := ch.QueueDeclare(QUEUE_NAME, true, false, false, false, nil)
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
				ContentType: "text/plain",
				Body:        []byte("benchmark pub-consume"),
			})
			<-msgs
		}
	}()
	b.Logf("Processed %d messages", b.N)
}

func Benchmark_ClassicQueue_1to1_2Connection(b *testing.B) {
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

	q, err := ch.QueueDeclare(QUEUE_NAME, true, false, false, false, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer ch.QueueDelete(q.Name, false, false, false)

	ch.QueuePurge(q.Name, false)

	msgs, err := ch.Consume(QUEUE_NAME, "", true, false, false, false, nil)
	if err != nil {
		b.Fatal(err)
	}

	go func() {
		conn, _ := amqp.Dial(RABBITMQ_URL)
		defer conn.Close()
		ch, _ := conn.Channel()
		defer ch.Close()
		for i := 0; i < b.N; i++ {
			_ = ch.Publish("", QUEUE_NAME, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("benchmark pub-consume"),
			})
		}
	}()

	for i := 0; i < b.N; i++ {
		<-msgs
	}
	b.Logf("Processed %d messages", b.N)
}

func Benchmark_ClassicQueue_1to1_1Connection(b *testing.B) {
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

	q, err := ch.QueueDeclare(QUEUE_NAME, true, false, false, false, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer ch.QueueDelete(q.Name, false, false, false)

	ch.QueuePurge(q.Name, false)

	msgs, err := ch.Consume(QUEUE_NAME, "", true, false, false, false, nil)
	if err != nil {
		b.Fatal(err)
	}

	go func() {
		ch, _ := conn.Channel()
		defer ch.Close()
		for i := 0; i < b.N; i++ {
			_ = ch.Publish("", QUEUE_NAME, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("benchmark pub-consume"),
			})
		}
	}()

	for i := 0; i < b.N; i++ {
		<-msgs
	}
	b.Logf("Processed %d messages", b.N)
}

func Benchmark_ClassicQueue_1toN_NConnection(b *testing.B) {
	n := 5
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

	q, err := ch.QueueDeclare(QUEUE_NAME, true, false, false, false, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer ch.QueueDelete(q.Name, false, false, false)

	ch.QueuePurge(q.Name, false)

	var consumed int64
	var mu sync.Mutex
	done := make(chan struct{})

	for i := 0; i < n; i++ {
		go func() {
			conn, _ := amqp.Dial(RABBITMQ_URL)
			defer conn.Close()
			ch, _ := conn.Channel()
			defer ch.Close()
			msgs, _ := ch.Consume(QUEUE_NAME, "", true, false, false, false, nil)
			for range msgs {
				mu.Lock()
				consumed++
				if consumed >= int64(b.N) {
					mu.Unlock()
					done <- struct{}{}
					return
				}
				mu.Unlock()
			}
		}()
	}

	go func() {
		ch, _ := conn.Channel()
		defer ch.Close()
		for i := 0; i < b.N; i++ {
			_ = ch.Publish("", QUEUE_NAME, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("benchmark pub-consume"),
			})
		}
	}()

	// 等待所有訊息被消費
	<-done
	b.Logf("Processed %d messages", b.N)
}

func Benchmark_ClassicQueue_1toN_1Channel(b *testing.B) {
	n := 5
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

	q, err := ch.QueueDeclare(QUEUE_NAME, true, false, false, false, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer ch.QueueDelete(q.Name, false, false, false)

	ch.QueuePurge(q.Name, false)

	var consumed int64
	var mu sync.Mutex
	done := make(chan struct{})

	for i := 0; i < n; i++ {
		go func() {
			ch, _ := conn.Channel()
			defer ch.Close()
			msgs, _ := ch.Consume(QUEUE_NAME, "", true, false, false, false, nil)
			for range msgs {
				mu.Lock()
				consumed++
				if consumed >= int64(b.N) {
					mu.Unlock()
					done <- struct{}{}
					return
				}
				mu.Unlock()
			}
		}()
	}

	go func() {
		ch, _ := conn.Channel()
		defer ch.Close()
		for i := 0; i < b.N; i++ {
			_ = ch.Publish("", QUEUE_NAME, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("benchmark pub-consume"),
			})
		}
	}()

	// 等待所有訊息被消費
	<-done
	b.Logf("Processed %d messages", b.N)
}

/*
❯ go test -bench=. -v  .
Alias tip: got -bench=. -v  .
goos: linux
goarch: amd64
pkg: rabbitmq/classic/tests
cpu: Intel(R) Core(TM) i5-10210U CPU @ 1.60GHz
Benchmark_ClassicQueue
    benchmark_test.go:54: Processed 1 messages
    benchmark_test.go:54: Processed 52 messages
    benchmark_test.go:54: Processed 2011 messages
    benchmark_test.go:54: Processed 4856 messages
Benchmark_ClassicQueue-4                            4856            240948 ns/op
Benchmark_ClassicQueue_1to1_2Connection
    benchmark_test.go:99: Processed 1 messages
    benchmark_test.go:99: Processed 69 messages
    benchmark_test.go:99: Processed 4290 messages
    benchmark_test.go:99: Processed 40593 messages
    benchmark_test.go:99: Processed 58627 messages
Benchmark_ClassicQueue_1to1_2Connection-4          58627             20706 ns/op
Benchmark_ClassicQueue_1to1_1Connection
    benchmark_test.go:142: Processed 1 messages
    benchmark_test.go:142: Processed 86 messages
    benchmark_test.go:142: Processed 6296 messages
    benchmark_test.go:142: Processed 45115 messages
    benchmark_test.go:142: Processed 54946 messages
Benchmark_ClassicQueue_1to1_1Connection-4          54946             22270 ns/op
Benchmark_ClassicQueue_1toN_NConnection
    benchmark_test.go:204: Processed 1 messages
    benchmark_test.go:204: Processed 60 messages
    benchmark_test.go:204: Processed 2710 messages
    benchmark_test.go:204: Processed 34705 messages
    benchmark_test.go:204: Processed 51783 messages
Benchmark_ClassicQueue_1toN_NConnection-4          51783             24002 ns/op
Benchmark_ClassicQueue_1toN_1Channel
    benchmark_test.go:264: Processed 1 messages
    benchmark_test.go:264: Processed 70 messages
    benchmark_test.go:264: Processed 3802 messages
    benchmark_test.go:264: Processed 41850 messages
    benchmark_test.go:264: Processed 51993 messages
Benchmark_ClassicQueue_1toN_1Channel-4             51993             23430 ns/op
PASS
ok      rabbitmq/classic/tests  10.888s
  ~/workspace/go-testcode/rabbitmq/classic/tests   master ?1 ─────────────────────────────────────────────────────────────────────────────────────────────────────────────  11s root@08cc3663a066  12:35:24
❯ go test -bench=. -v  .
Alias tip: got -bench=. -v  .
goos: linux
goarch: amd64
pkg: rabbitmq/classic/tests
cpu: Intel(R) Core(TM) i5-10210U CPU @ 1.60GHz
Benchmark_ClassicQueue
    benchmark_test.go:54: Processed 1 messages
    benchmark_test.go:54: Processed 52 messages
    benchmark_test.go:54: Processed 2101 messages
    benchmark_test.go:54: Processed 4922 messages
Benchmark_ClassicQueue-4                            4922            243767 ns/op
Benchmark_ClassicQueue_1to1_2Connection
    benchmark_test.go:99: Processed 1 messages
    benchmark_test.go:99: Processed 66 messages
    benchmark_test.go:99: Processed 4226 messages
    benchmark_test.go:99: Processed 42728 messages
    benchmark_test.go:99: Processed 56430 messages
Benchmark_ClassicQueue_1to1_2Connection-4          56430             21134 ns/op
Benchmark_ClassicQueue_1to1_1Connection
    benchmark_test.go:142: Processed 1 messages
    benchmark_test.go:142: Processed 86 messages
    benchmark_test.go:142: Processed 6110 messages
    benchmark_test.go:142: Processed 46446 messages
Benchmark_ClassicQueue_1to1_1Connection-4          46446             22491 ns/op
Benchmark_ClassicQueue_1toN_NConnection
    benchmark_test.go:204: Processed 1 messages
    benchmark_test.go:204: Processed 60 messages
    benchmark_test.go:204: Processed 2595 messages
    benchmark_test.go:204: Processed 33852 messages
    benchmark_test.go:204: Processed 54078 messages
Benchmark_ClassicQueue_1toN_NConnection-4          54078             23352 ns/op
Benchmark_ClassicQueue_1toN_1Channel
    benchmark_test.go:264: Processed 1 messages
    benchmark_test.go:264: Processed 55 messages
    benchmark_test.go:264: Processed 2194 messages
    benchmark_test.go:264: Processed 32583 messages
    benchmark_test.go:264: Processed 51571 messages
Benchmark_ClassicQueue_1toN_1Channel-4             51571             22092 ns/op
PASS
ok      rabbitmq/classic/tests  9.484s
  ~/workspace/go-testcode/rabbitmq/classic/tests   master ?1 ─────────────────────────────────────────────────────────────────────────────────────────────────────────────  10s root@08cc3663a066  12:35:41
❯ go test -bench=. -v  .
Alias tip: got -bench=. -v  .
goos: linux
goarch: amd64
pkg: rabbitmq/classic/tests
cpu: Intel(R) Core(TM) i5-10210U CPU @ 1.60GHz
Benchmark_ClassicQueue
    benchmark_test.go:54: Processed 1 messages
    benchmark_test.go:54: Processed 50 messages
    benchmark_test.go:54: Processed 1840 messages
    benchmark_test.go:54: Processed 4866 messages
Benchmark_ClassicQueue-4                            4866            240391 ns/op
Benchmark_ClassicQueue_1to1_2Connection
    benchmark_test.go:99: Processed 1 messages
    benchmark_test.go:99: Processed 64 messages
    benchmark_test.go:99: Processed 4111 messages
    benchmark_test.go:99: Processed 36886 messages
    benchmark_test.go:99: Processed 55501 messages
Benchmark_ClassicQueue_1to1_2Connection-4          55501             20929 ns/op
Benchmark_ClassicQueue_1to1_1Connection
    benchmark_test.go:142: Processed 1 messages
    benchmark_test.go:142: Processed 79 messages
    benchmark_test.go:142: Processed 5613 messages
    benchmark_test.go:142: Processed 44415 messages
Benchmark_ClassicQueue_1to1_1Connection-4          44415             22807 ns/op
Benchmark_ClassicQueue_1toN_NConnection
    benchmark_test.go:204: Processed 1 messages
    benchmark_test.go:204: Processed 57 messages
    benchmark_test.go:204: Processed 2487 messages
    benchmark_test.go:204: Processed 34858 messages
    benchmark_test.go:204: Processed 52636 messages
Benchmark_ClassicQueue_1toN_NConnection-4          52636             23338 ns/op
Benchmark_ClassicQueue_1toN_1Channel
    benchmark_test.go:264: Processed 1 messages
    benchmark_test.go:264: Processed 67 messages
    benchmark_test.go:264: Processed 4434 messages
    benchmark_test.go:264: Processed 44216 messages
Benchmark_ClassicQueue_1toN_1Channel-4             44216             23227 ns/op
PASS
ok      rabbitmq/classic/tests  8.379s
  ~/workspace/go-testcode/rabbitmq/classic/tests   master ?1 ──────────────────────────────────────────────────────────────────────────────────────────────────────────────  9s root@08cc3663a066  12:35:56
❯
*/
