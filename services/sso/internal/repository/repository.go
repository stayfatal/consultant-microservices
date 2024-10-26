package repository

import (
	"cm/services/sso/internal/interfaces"
	"cm/services/sso/internal/models"

	"github.com/pkg/errors"
)

type postgresRepo struct {
	db interfaces.DB
}

func New(db interfaces.DB) interfaces.Repository {
	return &postgresRepo{db: db}
}

func (repo *postgresRepo) CreateUser(user models.User) (int, error) {
	query := `INSERT INTO users
	(name,email,password,is_consultant)
	VALUES (:name,:email,:password,:is_consultant)
	RETURNING id`
	id := -1
	rows, err := repo.db.NamedQuery(query, user)
	if err != nil {
		return -1, errors.Wrap(err, "calling sqlx NamedQuery")
	}
	rows.Next()
	err = rows.Scan(&id)
	rows.Close()
	return id, errors.Wrap(err, "scanning rows")
}

func (repo *postgresRepo) GetUserByEmail(user models.User) (models.User, error) {
	foundedUser := models.User{}
	err := repo.db.Get(&foundedUser, "SELECT * FROM users WHERE email = $1", user.Email)
	return foundedUser, errors.Wrap(err, "calling sqlx Get")
}
