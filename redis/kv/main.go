package main

import (
	"context"
	"flag"
	"fmt"

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

	key, value := "getset-key", "value"

	// 設定 key
	err = rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}

	// 檢查 key 是否存在
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	if exists == 1 {
		fmt.Println("Key exists")
	}

	// 取得 key
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("Key does not exists")
		} else {
			panic(err)
		}
	}
	fmt.Println("key:", val)

	// 刪除 key
	err = rdb.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
