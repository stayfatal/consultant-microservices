package repository

import (
	"cm/services/chat/internal/interfaces"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) interfaces.Repository {
	return &repository{db}
}

func (repo *repository) CreateConsultant() {
	repo.db.Exec("")
}
