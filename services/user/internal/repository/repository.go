package repository

import (
	"cm/services/user/internal/interfaces"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db sqlx.DB
}

func New(db sqlx.DB) interfaces.Repository {
	return &repository{db: db}
}

func (repo *repository) CreateUser() {

}

func (repo *repository) GetUserByName() {

}
