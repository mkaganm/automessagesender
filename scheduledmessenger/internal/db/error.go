package db

import "errors"

var (
	ErrDSNNotProvided     = errors.New("MySQL DSN not provided")
	ErrMySQLConnection    = errors.New("failed to connect to MySQL")
	ErrPingFailed         = errors.New("failed to ping MySQL server")
	ErrReconnectionFailed = errors.New("failed to reconnect to MySQL")
	ErrDBCloseFailed      = errors.New("failed to close MySQL connection")
)
