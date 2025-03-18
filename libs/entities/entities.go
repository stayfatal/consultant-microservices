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

type Chat struct {
	Id           int `json:"id" db:"id"`
	ConsultantId int `json:"consultant_id" db:"consultant_id"`
	UserId       int `json:"user_id" db:"user_id"`
}

type Message struct {
	Id      int    `json:"id" db:"id"`
	ChatId  int    `json:"chat_id" db:"chat_id"`
	UserId  int    `json:"user_id" db:"user_id"`
	Message string `json:"message" db:"message"`
}
