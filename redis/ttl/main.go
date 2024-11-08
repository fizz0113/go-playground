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

	key, value := "temp", "xxx-token..."
	expiration := 10 * time.Second

	// 設定一個 key 含有超時設定
	err = rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		panic(err)
	}

	// 等 1s
	time.Sleep(time.Second)

	// 取得超時 這邊輸出會剩 9s
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("TTL:", ttl)

	// 刪除 key
	err = rdb.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
