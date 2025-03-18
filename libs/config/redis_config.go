package config

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisDB() (*redis.Client, error) {
	viper.AutomaticEnv()
	params := map[string]string{
		"PASSWORD": "",
		"HOST":     "",
		"PORT":     "",
	}

	for key := range params {
		paramKey := fmt.Sprintf("REDIS_%s", key)
		err := viper.BindEnv(paramKey)
		if err != nil {
			return nil, errors.Wrap(err, "binding env var")
		}
		params[key] = viper.GetString(paramKey)
	}

	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", params["HOST"], params["PORT"]),
		Password: params["PASSWORD"],
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return db, errors.Wrap(db.Ping(ctx).Err(), "pinging redis")
}
