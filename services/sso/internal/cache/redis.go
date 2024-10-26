package cache

import (
	"cm/services/sso/internal/interfaces"
	"cm/services/sso/internal/models"
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	db *redis.Client
}

func New(db *redis.Client) interfaces.CacheDB {
	return &redisRepo{db: db}
}

func (rr *redisRepo) SetUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	binary, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = rr.db.Set(ctx, user.Email, binary, time.Minute*10).Err()
	return errors.Wrap(err, "setting into redis")
}

func (rr *redisRepo) GetUser(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	got := models.User{}
	result, err := rr.db.Get(ctx, user.Email).Result()
	if err != nil {
		return models.User{}, err
	}
	err = json.Unmarshal([]byte(result), &got)
	if err != nil {
		return models.User{}, err
	}
	return got, errors.Wrap(err, "getting from redis")
}
