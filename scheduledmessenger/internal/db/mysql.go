package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
)

type MySQLClient struct {
	DB *gorm.DB
}

var (
	mysqlClientInstance *MySQLClient
	once                sync.Once
)

// Initialize sets up MySQL connection and configures the connection pool
func Initialize() error {
	var err error
	once.Do(func() {
		err = connect()
	})
	return err
}

// connect establishes a new connection to MySQL
func connect() error {
	dsn := os.Getenv("MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Connection pool settings with default values
	sqlDB.SetMaxOpenConns(0)    // Default is unlimited (0)
	sqlDB.SetMaxIdleConns(2)    // Default is 2 idle connections
	sqlDB.SetConnMaxLifetime(0) // Default is unlimited lifetime (0)
	sqlDB.SetConnMaxIdleTime(0) // Default is unlimited idle time (0)

	log.Println("MySQL connection successful")
	mysqlClientInstance = &MySQLClient{DB: db}
	return nil
}

// Reconnect checks the connection and reconnects if needed
func (client *MySQLClient) Reconnect() error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return err
	}

	// Ping the database to check if the connection is alive
	if err := sqlDB.Ping(); err != nil {
		log.Println("Lost connection to MySQL. Reconnecting...")
		return connect() // Re-establish connection
	}

	return nil
}

// GetInstance returns the singleton instance of MySQLClient.
func GetInstance() *MySQLClient {
	return mysqlClientInstance
}

// Close terminates the database connection
func (client *MySQLClient) Close() error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Wrapped CRUD Functions with Reconnect Check

// Create inserts a new record into the database
func (client *MySQLClient) Create(value interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.DB.Create(value).Error
}

// Find retrieves records from the database based on conditions
func (client *MySQLClient) Find(value interface{}, conditions ...interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.DB.Find(value, conditions...).Error
}

// First retrieves the first matching record based on conditions
func (client *MySQLClient) First(value interface{}, conditions ...interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.DB.First(value, conditions...).Error
}

// Update updates an existing record in the database
func (client *MySQLClient) Update(value interface{}) error {
	if err := client.Reconnect(); err != nil {
		return err
	}
	return client.DB.Save(value).Error
}

// Delete removes a record from the database based on conditions
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

// RawQuery executes a raw SQL query on the database
func (client *MySQLClient) RawQuery(query string, values ...interface{}) (*gorm.DB, error) {
	if err := client.Reconnect(); err != nil {
		return nil, err
	}
	result := client.DB.Raw(query, values...)
	return result, result.Error
}

// Query is a general-purpose query function that allows chaining conditions
func (client *MySQLClient) Query() *gorm.DB {
	if err := client.Reconnect(); err != nil {
		log.Printf("Failed to reconnect to MySQL: %v", err)
		return nil
	}
	return client.DB
}
