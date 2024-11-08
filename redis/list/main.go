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

	key, values := "mylist", []interface{}{"a", "b", "c", "d"}

	// 設定 mylist: ["a", "b", "c", "d"]
	err = rdb.RPush(ctx, key, values...).Err()
	if err != nil {
		panic(err)
	}

	// 左推也行 mylist: ["left", "a", "b", "c", "d"]
	err = rdb.LPush(ctx, key, "left").Err()
	if err != nil {
		panic(err)
	}

	// 取出 mylist
	vals, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("mylist:", vals)

	// 刪除 mylist
	err = rdb.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
