package models

import "time"

type User struct {
	Id        int       `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
