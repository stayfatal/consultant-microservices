package service

import (
	"cm/services/sso/internal/auth"
	"cm/services/sso/internal/interfaces"
	"cm/services/sso/internal/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"
)

type service struct {
	repo interfaces.Repository
}

func New(repo interfaces.Repository) interfaces.Service {
	return &service{repo: repo}
}

func (svc *service) Register(user models.User) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "generating hashed password")
	}
	user.Password = string(hashedPass)

	id, err := svc.repo.CreateUser(user)
	if err != nil {
		return "", errors.Wrap(err, "creating user service level")
	}
	user.Id = id

	token, err := auth.CreateToken(user.Id)
	if err != nil {
		return "", errors.Wrap(err, "creating token")
	}

	return token, nil
}

func (svc *service) Login(user models.User) (string, error) {
	foundUser, err := svc.repo.GetUserByEmail(user)
	if err != nil {
		return "", errors.Wrap(err, "getting user by email service level")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		return "", errors.Wrap(err, "comparing hash and password")
	}

	token, err := auth.CreateToken(foundUser.Id)
	if err != nil {
		return "", errors.Wrap(err, "creating token")
	}

	return token, nil
}
