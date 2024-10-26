package cache

import (
	"cm/services/sso/internal/dbtest"
	"cm/services/sso/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGetUser(t *testing.T) {
	db, err := dbtest.PrepareCachingDB()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

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
