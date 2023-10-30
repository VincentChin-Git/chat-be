package storage

import (
	"client-server-go/config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func ConnectRedis() {
	configGet := config.GetConfig()

	redisClient = redis.NewClient(&redis.Options{
		Addr:     configGet.RedisAddr,
		Password: configGet.RedisPass,
		DB:       0, // use default DB
	})
}

func WriteRedis(key string, value string, expired time.Duration) bool {
	ctx := context.Background()
	err := redisClient.Set(ctx, key, value, expired).Err()

	return err != nil
}

func ReadRedis(key string) interface{} {
	ctx := context.Background()
	val, err := redisClient.Get(ctx, key).Result()

	if err != nil {
		return false
	}

	return val
}
