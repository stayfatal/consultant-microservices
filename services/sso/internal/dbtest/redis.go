package dbtest

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func PrepareCachingDB() (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:8081",
		Password: "",
		DB:       1,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := db.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return db, nil
}
