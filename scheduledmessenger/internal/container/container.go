package container

import (
	"log"
	"scheduledmessenger/internal/db"
	"scheduledmessenger/internal/message"
	"scheduledmessenger/internal/queue"
	"scheduledmessenger/internal/scheduler"
	"scheduledmessenger/pkg/logger"
)

type Container struct {
	Logger            *logger.Logger
	MySQLClient       *db.MySQLClient
	RabbitMQClient    *queue.RabbitMQClient
	Scheduler         *scheduler.CronClient
	MessageRepository *message.MessageRepository
	MessageService    *message.Service
}

// Create initializes the infrastructure components
func Create() *Container {
	container := new(Container)

	logger := logger.NewLogger()
	container.Logger = logger

	// Initialize MySQL client
	mysqlClient, err := db.CreateClient()
	if err != nil {
		log.Fatalf("Failed to create MySQL client: %v", err)
	}
	container.MySQLClient = mysqlClient

	// Initialize RabbitMQ client
	rabbitMQClient, err := queue.CreateQueueClient()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	container.RabbitMQClient = rabbitMQClient

	scheduler := scheduler.NewCronJobService()
	container.Scheduler = scheduler

	return container
}

// Initialize sets up the services
func (c *Container) Initialize() {
	// Initialize message repository
	c.MessageRepository = message.CreateRepository(c)

	// Initialize message service
	c.MessageService = message.CreateService(c)
}
