package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLClient struct {
	db *gorm.DB
}

func CreateClient() (*MySQLClient, error) {
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

	return &MySQLClient{db: db}, nil
}

// Reconnect checks the connection and reconnects if needed
func (client *MySQLClient) Reconnect() error {
	sqlDB, err := client.db.DB()
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
		client.db = db
	}

	return nil
}

// Close terminates the database connection
func (client *MySQLClient) Close() error {
	sqlDB, err := client.db.DB()
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
	return client.db.Create(value).Error
}

func (client *MySQLClient) Find(value interface{}, conditions ...interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.db.Find(value, conditions...).Error
}

func (client *MySQLClient) First(value interface{}, conditions ...interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.db.First(value, conditions...).Error
}

func (client *MySQLClient) Update(value interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.db.Save(value).Error
}

func (client *MySQLClient) Delete(value interface{}, conditions ...interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	query := client.db
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Delete(value).Error
}

func (client *MySQLClient) RawQuery(query string, values ...interface{}) (*gorm.DB, error) {
	if err := client.Reconnect(); err != nil {
		return nil, err
	}
	result := client.db.Raw(query, values...)
	return result, result.Error
}

func (client *MySQLClient) Query() *gorm.DB {
	if err := client.Reconnect(); err != nil {
		log.Printf("Failed to reconnect to MySQL: %v", err)
		return nil
	}
	return client.db
}
