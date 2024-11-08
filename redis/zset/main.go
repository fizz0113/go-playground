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

	key, zs := "scores", []redis.Z{{Score: 10, Member: "ian10"}, {Score: 1, Member: "ian1"}}

	// 設定 sorted set
	err = rdb.ZAdd(ctx, key, zs...).Err()
	if err != nil {
		panic(err)
	}

	// 取出 sorted set, 他會按照 Score 排序好
	res, err := rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	for _, z := range res {
		fmt.Println(z.Member, z.Score)
	}

	// 刪除 scores
	err = rdb.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
