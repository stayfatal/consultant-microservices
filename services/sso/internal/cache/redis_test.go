package cache

import (
	"cm/services/sso/internal/models"
	"cm/services/sso/internal/testhelpers"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGetUser(t *testing.T) {
	ctx := context.Background()
	container, db, err := testhelpers.ConfigureRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupRedisContainer(t, container, db)
	cache := New(db)

	user := models.User{
		Name:         "test",
		Email:        "test@testmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	err = cache.SetUser(user)
	if err != nil {
		t.Fatal(err)
	}

	foundUser, err := cache.GetUser(user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, user, foundUser)
}
