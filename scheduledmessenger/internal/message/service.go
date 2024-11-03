package message

import (
	"encoding/json"
	"scheduledmessenger/internal/container"
	"scheduledmessenger/internal/queue"
	"scheduledmessenger/internal/scheduler"
)

const queueName = "MESSAGE_SEND_QUEUE"

type Service struct {
	repository  *MessageRepository
	queueClient *queue.RabbitMQClient
	scheduler   *scheduler.CronClient
}

func CreateService(c *container.Container) *Service {
	return &Service{
		repository:  c.MessageRepository,
		queueClient: c.RabbitMQClient,
		scheduler:   c.Scheduler,
	}
}

func (s *Service) Run() {
	s.scheduler.AddJob("*/2 * * * * *", func() {
		err := s.sendMessage()

		if err != nil {
			// fixme :
		}
	})

	s.scheduler.Start()

	// Keep the program running
	select {}
}

func (s *Service) sendMessage() error {
	unsentTwoMessages, err := s.repository.GetUnsentTwoMessages()

	if err != nil {
		return err
	}

	for _, message := range unsentTwoMessages {

		s.repository.MarkAsProcessed(&message)

		msgContext, err := json.Marshal(message)

		if err != nil {
			return err
		}

		err = s.queueClient.PublishMessage(msgContext, queueName)

		if err != nil {
			return err
		}
	}

	return nil
}
