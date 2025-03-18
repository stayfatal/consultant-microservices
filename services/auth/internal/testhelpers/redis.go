package testhelpers

import (
	"cm/libs/config"
	"testing"

	"github.com/redis/go-redis/v9"
)

func PrepareRedis(t *testing.T) (*redis.Client, error) {
	db, err := config.NewRedisDB()
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
