package cache

import (
	"context"
	"fmt"

	"food-serve.com/pkg/config"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RedisClient(cfg config.Config) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	defer rdb.Close()

	pong, err := rdb.Ping(ctx).Result()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Redis:", pong)
}
