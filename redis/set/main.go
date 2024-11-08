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

	key, members := "tags", []interface{}{"go", "redis", "database"}

	// 設定 set, 概念跟 map[string]struct{} 相同, 無序/去重複
	err = rdb.SAdd(ctx, key, members...).Err()
	if err != nil {
		panic(err)
	}

	// 取出 set
	tags, err := rdb.SMembers(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("tags:", tags)

	// 刪除 tags
	err = rdb.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
