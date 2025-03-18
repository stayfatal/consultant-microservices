package service

import (
	"cm/libs/entities"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type repositoryMock struct {
	counter  int
	expected entities.User
}

func (repo *repositoryMock) CreateUser(user entities.User) (int, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(repo.expected.Password))
	if err != nil {
		return -1, err
	}
	repo.counter++
	return repo.expected.Id, nil
}

func (repo *repositoryMock) GetUserByEmail(user entities.User) (entities.User, error) {
	retUser := user
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(retUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, err
	}
	retUser.Password = string(hashedPass)
	repo.counter++
	return retUser, nil
}

type cacheMock struct {
	counter  int
	expected entities.User
	noUser   bool
}

func (repo *cacheMock) SetUser(user entities.User) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(repo.expected.Password))
	if err != nil {
		return err
	}
	repo.counter++
	return nil
}

func (repo *cacheMock) GetUser(user entities.User) (entities.User, error) {
	if repo.noUser {
		return entities.User{}, redis.Nil
	}
	retUser := user
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(retUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, err
	}
	retUser.Password = string(hashedPass)
	repo.counter++
	return retUser, nil
}
