package entities

import (
	"time"
)

type User struct {
	Id           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"password" db:"password"`
	IsConsultant bool      `json:"is_consultant" db:"is_consultant"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
