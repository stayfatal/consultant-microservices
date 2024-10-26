package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisDb(cfg *Config) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_ADDR,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := db.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return db, nil
}
