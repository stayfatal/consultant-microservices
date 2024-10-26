package cache

import (
	"cm/services/sso/internal/models"
	"cm/services/sso/internal/testhelpers"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetUser(t *testing.T) {
	ctx := context.Background()
	container, db, err := testhelpers.ConfigureRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupRedisContainer(t, container, db)
	cache := New(db)

	expected := models.User{
		Name:         "test",
		Email:        "test@testmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	err = cache.SetUser(expected)
	if err != nil {
		t.Fatal(err)
	}

	got := models.User{}
	result, err := db.Get(ctx, expected.Email).Result()
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal([]byte(result), &got)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, got)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	container, db, err := testhelpers.ConfigureRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupRedisContainer(t, container, db)
	cache := New(db)

	expected := models.User{
		Name:         "test",
		Email:        "test@testmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	binary, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Set(ctx, expected.Email, binary, time.Minute*10).Err()
	if err != nil {
		t.Fatal()
	}

	got, err := cache.GetUser(expected)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, got)
}
