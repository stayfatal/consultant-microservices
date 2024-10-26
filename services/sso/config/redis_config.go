package config

import (
	"cm/services/sso/utils"
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	viperr "github.com/spf13/viper"
)

type RedisConfig struct {
	IMAGE    string
	HOST     string
	PORT     int
	PASSWORD string
}

func LoadRedisConfig() (*RedisConfig, error) {
	path, err := utils.GetPath("services/sso/config/redis_config.env")
	if err != nil {
		return nil, err
	}
	viperr.SetConfigFile(path)
	err = viperr.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	cfg := &RedisConfig{
		IMAGE:    viperr.GetString("IMAGE"),
		HOST:     viperr.GetString("HOST"),
		PORT:     viperr.GetInt("PORT"),
		PASSWORD: viperr.GetString("PASSWORD"),
	}

	return cfg, nil
}

func NewRedisDb(cfg *RedisConfig) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.HOST, cfg.PORT),
		Password: cfg.PASSWORD,
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
