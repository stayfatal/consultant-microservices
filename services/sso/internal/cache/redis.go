package cache

import (
	"cm/internal/entities"
	"cm/services/sso/internal/interfaces"
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

func (rr *redisRepo) SetUser(user entities.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	binary, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = rr.db.Set(ctx, user.Email, binary, time.Minute*10).Err()
	return errors.Wrap(err, "setting into redis")
}

func (rr *redisRepo) GetUser(user entities.User) (entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	got := entities.User{}
	result, err := rr.db.Get(ctx, user.Email).Result()
	if err != nil {
		return entities.User{}, err
	}
	err = json.Unmarshal([]byte(result), &got)
	if err != nil {
		return entities.User{}, err
	}
	return got, errors.Wrap(err, "getting from redis")
}
