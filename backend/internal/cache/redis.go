package cache

import (
	"context"
	"fmt"
	"time"

	"food-serve.com/pkg/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

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

func (c Cache) GetKey(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c Cache) SetKey(ctx context.Context, key string, value interface{}, time time.Duration) error {
	return c.client.Set(ctx, key, value, time).Err()
}
