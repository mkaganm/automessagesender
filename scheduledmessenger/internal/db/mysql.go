package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLClient struct {
	DB *gorm.DB
}

var (
	ErrDSNNotProvided     = errors.New("MySQL DSN not provided")
	ErrMySQLConnection    = errors.New("failed to connect to MySQL")
	ErrPingFailed         = errors.New("failed to ping MySQL server")
	ErrReconnectionFailed = errors.New("failed to reconnect to MySQL")
	ErrDBCloseFailed      = errors.New("failed to close MySQL connection")
)

// NewMySQLClient initializes a new MySQLClient instance.
// It gets the DSN directly from the environment variable.
func NewMySQLClient() (*MySQLClient, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return nil, ErrDSNNotProvided
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMySQLConnection, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMySQLConnection, err)
	}

	// Set connection pool properties
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	log.Println("MySQL connection established successfully")

	return &MySQLClient{DB: db}, nil
}

// Reconnect checks the connection and reconnects if needed
func (client *MySQLClient) Reconnect() error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPingFailed, err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Println("Lost connection to MySQL. Reconnecting...")
		dsn := os.Getenv("MYSQL_DSN")
		if dsn == "" {
			return ErrDSNNotProvided
		}

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("%w: %v", ErrReconnectionFailed, err)
		}
		client.DB = db
	}

	return nil
}

// Close terminates the database connection
func (client *MySQLClient) Close() error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBCloseFailed, err)
	}
	return sqlDB.Close()
}

// CRUD methods with Reconnect check

func (client *MySQLClient) Create(value interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.DB.Create(value).Error
}

func (client *MySQLClient) Find(value interface{}, conditions ...interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.DB.Find(value, conditions...).Error
}

func (client *MySQLClient) First(value interface{}, conditions ...interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.DB.First(value, conditions...).Error
}

func (client *MySQLClient) Update(value interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.DB.Save(value).Error
}

func (client *MySQLClient) Delete(value interface{}, conditions ...interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	query := client.DB
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Delete(value).Error
}

func (client *MySQLClient) RawQuery(query string, values ...interface{}) (*gorm.DB, error) {
	if err := client.Reconnect(); err != nil {
		return nil, err
	}
	result := client.DB.Raw(query, values...)
	return result, result.Error
}

func (client *MySQLClient) Query() *gorm.DB {
	if err := client.Reconnect(); err != nil {
		log.Printf("Failed to reconnect to MySQL: %v", err)
		return nil
	}
	return client.DB
}
