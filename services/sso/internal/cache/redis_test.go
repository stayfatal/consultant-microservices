package cache

import (
	"cm/internal/entities"
	"cm/services/sso/internal/testhelpers"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSetUser(t *testing.T) {
	ctx := context.Background()
	db, err := testhelpers.PrepareRedis(t)
	if err != nil {
		t.Fatal(err)
	}

	cache := New(db)

	expected := entities.User{
		Name:         "test",
		Email:        fmt.Sprintf("test%s@gmail.com", uuid.New().String()),
		Password:     "123",
		IsConsultant: false,
	}

	err = cache.SetUser(expected)
	if err != nil {
		t.Fatal(err)
	}

	got := entities.User{}
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
	db, err := testhelpers.PrepareRedis(t)
	if err != nil {
		t.Fatal(err)
	}

	cache := New(db)

	expected := entities.User{
		Name:         "test",
		Email:        fmt.Sprintf("test%s@gmail.com", uuid.New().String()),
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
