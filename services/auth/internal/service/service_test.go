package service

import (
	"cm/libs/entities"
	"cm/libs/publicauth"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var expected = entities.User{
	Id:           1,
	Name:         "test",
	Email:        "test@testmail.com",
	Password:     "123",
	IsConsultant: false,
	CreatedAt:    time.Now(),
}

func TestRegister(t *testing.T) {
	cache := &cacheMock{expected: expected}
	repo := &repositoryMock{expected: expected}

	svc := New(repo, cache)

	token, err := svc.Register(expected)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := publicauth.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Id, claims.Id)
	assert.Equal(t, expected.Email, claims.Email)

	assert.Equal(t, 1, cache.counter)
	assert.Equal(t, 1, repo.counter)
}

func TestLogin(t *testing.T) {
	type testCase struct {
		user     entities.User
		repoSave bool
	}

	testCases := []testCase{
		{
			expected,
			false,
		},
		{
			expected,
			true,
		},
	}

	for _, test := range testCases {
		cache := &cacheMock{expected: expected}
		if test.repoSave {
			cache.noUser = true
		}
		repo := &repositoryMock{expected: expected}

		svc := New(repo, cache)

		token, err := svc.Login(expected)
		if err != nil {
			t.Fatal(err)
		}

		claims, err := publicauth.ValidateToken(token)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected.Id, claims.Id)
		assert.Equal(t, expected.Email, claims.Email)

		if test.repoSave {
			assert.Equal(t, 1, repo.counter)
		} else {
			assert.Equal(t, 1, cache.counter)
		}
	}
}
