package message

import (
	"scheduledmessenger/internal/container"
	"scheduledmessenger/internal/db"
)

// MessageRepository represents the message repository with a shared MySQLClient instance
type MessageRepository struct {
	client *db.MySQLClient
}

func CreateRepository(c *container.Container) *MessageRepository {
	return &MessageRepository{client: c.MySQLClient}
}

// GetUnsentMessages retrieves messages that have not been sent
func (r *MessageRepository) GetUnsentMessages() ([]Message, error) {
	var messages []Message
	err := r.client.Find(&messages, "is_send = ?", false)
	return messages, err
}

// GetUnsentTwoMessages retrieves the first 2 unsent messages, ordered by created_at
func (r *MessageRepository) GetUnsentTwoMessages() ([]Message, error) {
	var messages []Message
	err := r.client.Query().
		Where("is_send = ?", false).
		Order("created_at ASC").
		Limit(2).
		Find(&messages).Error
	return messages, err
}

// MarkAsSent updates a message to mark it as sent
func (r *MessageRepository) MarkAsSent(msg *Message) error {
	msg.IsSend = true
	return r.client.Update(msg)
}

// MarkAsProcessed updates a message to mark it as processed
func (r *MessageRepository) MarkAsProcessed(msg *Message) error {
	msg.IsSend = true
	return r.client.Update(msg)
}
