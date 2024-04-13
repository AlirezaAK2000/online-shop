package initializers

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisCtx = context.Background()
var RedisDB *redis.Client

func InitilizeRedisConnection() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}
