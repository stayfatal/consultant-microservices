package testhelpers

import (
	"cm/services/sso/config"
	"testing"

	"github.com/redis/go-redis/v9"
)

func PrepareRedis(t *testing.T) (*redis.Client, error) {
	cfg, err := config.LoadRedisConfig()
	if err != nil {
		return nil, err
	}

	db, err := config.NewRedisDb(cfg)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	})

	return db, nil
}
