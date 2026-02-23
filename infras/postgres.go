package infras

import (
	"fmt"
	"time"

	"github.com/ariashabry/boilerplate-go/helpers/env"
	applog "github.com/ariashabry/boilerplate-go/helpers/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	postgresMaxIdleConnection = 10
	postgresMaxOpenConnection = 10
	postgresConnMaxLifetime   = time.Hour
)

// PostgresConn holds both Read and Write database connections
type PostgresConn struct {
	Read  *gorm.DB
	Write *gorm.DB
}

// ProvidePostgresConn is the provider function for Wire DI
func ProvidePostgresConn(config *env.Config, log *applog.AppLog) *PostgresConn {
	return &PostgresConn{
		Read:  CreatePostgresReadConn(config, log),
		Write: CreatePostgresWriteConn(config, log),
	}
}

// CreatePostgresWriteConn creates a database connection for write access
func CreatePostgresWriteConn(config *env.Config, log *applog.AppLog) *gorm.DB {
	return CreatePostgresConnection(
		"write",
		config.DBPostgresWriteUser,
		config.DBPostgresWritePassword,
		config.DBPostgresWriteHost,
		config.DBPostgresWritePort,
		config.DBPostgresWriteName,
		config.DBPostgresWriteTimezone,
		config.DBPostgresMaxRetry,
		config.DBPostgresRetryWaitTime,
		log,
	)
}

// CreatePostgresReadConn creates a database connection for read access
func CreatePostgresReadConn(config *env.Config, log *applog.AppLog) *gorm.DB {
	return CreatePostgresConnection(
		"read",
		config.DBPostgresReadUser,
		config.DBPostgresReadPassword,
		config.DBPostgresReadHost,
		config.DBPostgresReadPort,
		config.DBPostgresReadName,
		config.DBPostgresReadTimezone,
		config.DBPostgresMaxRetry,
		config.DBPostgresRetryWaitTime,
		log,
	)
}

// CreatePostgresConnection creates a database connection with retry logic
func CreatePostgresConnection(name, username, password, host, port, dbName, timeZone string, maxRetry, waitTime int, log *applog.AppLog) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, username, password, dbName, port, timeZone)

	var db *gorm.DB
	var err error

	for i := 0; i < maxRetry; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.WithFields(map[string]interface{}{
				"name":   name,
				"host":   host,
				"port":   port,
				"dbName": dbName,
			}).Info("Connected to Postgres database")

			// Configure connection pool
			sqlDB, err := db.DB()
			if err != nil {
				log.WithError(err).Error("Failed to get underlying SQL DB")
			} else {
				sqlDB.SetMaxIdleConns(postgresMaxIdleConnection)
				sqlDB.SetMaxOpenConns(postgresMaxOpenConnection)
				sqlDB.SetConnMaxLifetime(postgresConnMaxLifetime)
			}

			return db
		}

		log.WithFields(map[string]interface{}{
			"name":    name,
			"host":    host,
			"port":    port,
			"dbName":  dbName,
			"attempt": i + 1,
			"error":   err.Error(),
		}).Warn("Failed connecting to Postgres database, retrying")

		time.Sleep(time.Duration(waitTime) * time.Second)
	}

	log.WithFields(map[string]interface{}{
		"name":   name,
		"host":   host,
		"port":   port,
		"dbName": dbName,
		"error":  err.Error(),
	}).Fatal("Failed connecting to Postgres database after retries")

	return nil
}
