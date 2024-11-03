package queue

import "errors"

var (
	ErrEnvVarNotSet = errors.New("RABBITMQ_URL environment variable is not set")
	ErrConnFailed   = errors.New("failed to connect to RabbitMQ")
	ErrChanFailed   = errors.New("failed to open a channel")
)
