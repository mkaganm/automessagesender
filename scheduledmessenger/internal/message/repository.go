package message

import (
	"go.uber.org/dig"
	"log"
	"scheduledmessenger/internal/db"
)

// Repository represents the message repository with a shared MySQLClient instance
type Repository struct {
	client *db.MySQLClient
}

type RepositoryParams struct {
	dig.In
	DBClient *db.MySQLClient
}

var (
	repoInstance *Repository
)

// CreateRepository sets up the singleton instance of Repository
func CreateRepository(p RepositoryParams) *Repository {
	return &Repository{client: p.DBClient}
}

// GetRepository returns the singleton instance of Repository
func GetRepository() *Repository {
	if repoInstance == nil {
		log.Fatal("Repository not initialized. Call CreateRepository() first.")
	}
	return repoInstance
}

// GetUnsentMessages retrieves messages that have not been sent
func (r *Repository) GetUnsentMessages() ([]Message, error) {
	var messages []Message
	err := r.client.Find(&messages, "is_send = ?", false)
	return messages, err
}

// GetUnsentTwoMessages retrieves the first 2 unsent messages, ordered by created_at
func (r *Repository) GetUnsentTwoMessages() ([]Message, error) {
	var messages []Message
	err := r.client.Query().
		Where("is_send = ?", false).
		Order("created_at ASC").
		Limit(2).
		Find(&messages).Error
	return messages, err
}

// MarkAsSent updates a message to mark it as sent
func (r *Repository) MarkAsSent(msg *Message) error {
	msg.IsSend = true
	return r.client.Update(msg)
}
