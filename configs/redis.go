package configs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	Redis     *redis.Client
	onceRedis sync.Once
)

func ConnectToRedis() {
	onceRedis.Do(func() {
		redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		rdb := redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			DB:   redisDB,
		})

		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatalf("Failed to connect to redis: %v", err)
		}
<<<<<<< HEAD

=======
		//Redis = gs
>>>>>>> 5d459e3b2242427f8fb133f65a714d6ae77c113a
		Redis = rdb
	})
}
