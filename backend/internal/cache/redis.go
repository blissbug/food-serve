package rs

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RedisClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
