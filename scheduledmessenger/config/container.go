package config

import (
	"go.uber.org/dig"
	"scheduledmessenger/internal/db"
	"scheduledmessenger/internal/message"
	"scheduledmessenger/internal/queue"
)

func SetupContainer() (*dig.Container, error) {
	container := dig.New()

	providers := []interface{}{
		db.NewMySQLClient,
		queue.NewRabbitMQClient,
		message.CreateRepository,
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, err
		}
	}

	return container, nil
}
