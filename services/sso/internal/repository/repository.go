package repository

import (
	"cm/services/sso/internal/interfaces"
	"cm/services/sso/internal/models"

	"github.com/pkg/errors"
)

type repository struct {
	db interfaces.DB
}

func New(db interfaces.DB) interfaces.Repository {
	return &repository{db: db}
}

func (repo *repository) CreateUser(user models.User) (int, error) {
	query := `INSERT INTO users
	(name,email,password,is_consultant)
	VALUES (:name,:email,:password,:is_consultant)
	RETURNING id`
	id := -1
	rows, err := repo.db.NamedQuery(query, user)
	if err != nil {
		return -1, errors.Wrap(err, "creating user repository level")
	}
	rows.Next()
	err = rows.Scan(&id)
	return id, errors.Wrap(err, "creating user repository level")
}

func (repo *repository) GetUserByEmail(user models.User) (models.User, error) {
	foundedUser := models.User{}
	err := repo.db.Get(&foundedUser, "SELECT * FROM users WHERE email = $1", user.Email)
	return foundedUser, errors.Wrap(err, "getting user by email repository level")
}
