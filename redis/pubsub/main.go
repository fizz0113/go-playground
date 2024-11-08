package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisHostName = flag.String("redisHostName", "redis", "Specify redis host name")
	redisHostPort = flag.Uint("redisHostPort", 6379, "Specify redis host port")
	redisPassord  = flag.String("redisPassword", "mypassword", "Specify redis password")
	redisDB       = flag.Int("redisDB", 0, "Specify redis database")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", *redisHostName, *redisHostPort),
		Password: *redisPassord,
		DB:       *redisDB,
	})

	// 確認 redis 的連線
	err := rdb.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	topic := "topic"

	// 發送 event
	go func() {
		for {
			msg := fmt.Sprintf("current time: %s", time.Now())
			err := rdb.Publish(ctx, topic, msg).Err()
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second)
		}
	}()

	// 接收 event
	sub := rdb.Subscribe(ctx, topic)
	defer sub.Close()
	ch := sub.Channel()
	for msg := range ch {
		fmt.Printf("[%s] %s\n", msg.Channel, msg.Payload)
	}
}
