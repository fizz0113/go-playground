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

	key, hash := "myhash", map[string]interface{}{
		"name":  "ian",
		"email": "ian@example.com",
	}

	// 設定 myhash
	err = rdb.HSet(ctx, key, hash).Err()
	if err != nil {
		panic(err)
	}

	// 取出 name 欄位
	name, err := rdb.HGet(ctx, key, "name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("name:", name)

	// 取出 email 欄位
	email, err := rdb.HGet(ctx, key, "email").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("email:", email)

	// 更新 email 欄位
	m := map[string]interface{}{"email": "updated@example.com"}
	err = rdb.HSet(ctx, key, m).Err()
	if err != nil {
		panic(err)
	}

	// 取出 email 欄位
	email, err = rdb.HGet(ctx, key, "email").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("updated email:", email)

	// 刪除 myhash
	err = rdb.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
