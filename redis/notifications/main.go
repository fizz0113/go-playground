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

	// 開啟設定
	// 事件類型選項：
	//   A - All (等同啟用全部事件類型)
	//   K - Keyspace events (例如 __keyspace@0__:mykey) → 表示 key 被誰操作
	//   E - Keyevent events  (例如 __keyevent@0__:set)   → 表示發生了什麼事件
	//
	// 命令類型事件：
	//   g - Generic commands (set, del, expire, rename 等通用命令)
	//   x - Expired events (key TTL 到期被刪除)
	//   e - Evicted events (key 因記憶體不足被淘汰)
	//
	// 特定資料結構命令事件：
	//   z - Sorted Set (zadd, zrem, zincrby)
	//   s - Set (sadd, srem)
	//   l - List (lpush, rpop, lpop)
	//   h - Hash (hset, hdel, hincrby)
	//   t - Stream (xadd, xdel, xtrim)
	//   d - Module commands (Redis module 執行的命令，例如 RedisJSON)
	//
	// 進階功能：
	//   m - Key miss events（GET 不存在 key）
	//   $ - String commands（setrange, append, incr 等）【Redis 7+】
	//   c - Touch commands（touch, getdel, getex）【Redis 6+】
	//   u - Stream consumer group commands（xreadgroup, xack 等）【Redis 7+】
	//   r - Read-only commands（get, smembers 等）【Redis 7.2+】
	err = rdb.ConfigSet(ctx, "notify-keyspace-events", "AKEgxl").Err()
	if err != nil {
		panic(err)
	}

	topic := "__key*@0__:*"

	// 接收 redis notifications
	sub := rdb.PSubscribe(ctx, topic)
	defer sub.Close()
	ch := sub.Channel()
	for msg := range ch {
		fmt.Printf("[%s] %s\n", msg.Channel, msg.Payload)
	}
}
