package repository

import (
	"cm/internal/entities"
	"cm/services/chat/internal/interfaces"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) interfaces.Repository {
	return &repository{db}
}

func (repo *repository) CreateChat(chat entities.Chat) (id int, err error) {
	err = repo.db.QueryRow("INSERT INTO chats (consultant_id,user_id) VALUES ($1,$2) RETURNING id", chat.ConsultantId, chat.UserId).Scan(&id)
	return
}

func (repo *repository) SaveMessage(message entities.Message) error {
	_, err := repo.db.Exec("INSERT INTO messages (chat_id,user_id,message) VALUES ($1,$2,$3)", message.ChatId, message.UserId, message.Message)
	return err
}
