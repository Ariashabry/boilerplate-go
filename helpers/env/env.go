package env

import (
	"sync"
	"time"

	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/spf13/viper"
)

type Config struct {
	// ======================
	// APP
	// ======================
	AppName     string `mapstructure:"APP_NAME"`
	AppRevision string `mapstructure:"APP_REVISION"`
	AppURL      string `mapstructure:"APP_URL"`
	AppHost     string `mapstructure:"APP_HOST"`
	AppPort     int    `mapstructure:"APP_PORT"`

	// ======================
	// CORS
	// ======================
	AppCorsAllowCredentials bool     `mapstructure:"APP_CORS_ALLOW_CREDENTIALS"`
	AppCorsAllowWildcard    bool     `mapstructure:"APP_CORS_ALLOW_WILDCARD"`
	AppCorsAllowedHeaders   []string `mapstructure:"APP_CORS_ALLOWED_HEADERS"`
	AppCorsAllowedMethods   []string `mapstructure:"APP_CORS_ALLOWED_METHODS"`
	AppCorsAllowedOrigins   []string `mapstructure:"APP_CORS_ALLOWED_ORIGINS"`
	AppCorsEnable           bool     `mapstructure:"APP_CORS_ENABLE"`
	AppCorsMaxAgeSeconds    int      `mapstructure:"APP_CORS_MAX_AGE_SECONDS"`

	// ======================
	// CACHE REDIS
	// ======================
	CacheRedisPrimaryHost     string        `mapstructure:"CACHE_REDIS_PRIMARY_HOST"`
	CacheRedisPrimaryPort     string        `mapstructure:"CACHE_REDIS_PRIMARY_PORT"`
	CacheRedisPrimaryPassword string        `mapstructure:"CACHE_REDIS_PRIMARY_PASSWORD"`
	CacheRedisPrimaryDB       int           `mapstructure:"CACHE_REDIS_PRIMARY_DB"`
	CacheDefaultExpiresIn     time.Duration `mapstructure:"CACHE_DEFAULT_EXPIRES_IN"`

	// ======================
	// MYSQL READ
	// ======================
	DBMySQLReadHost     string `mapstructure:"DB_MYSQL_READ_HOST"`
	DBMySQLReadPort     string `mapstructure:"DB_MYSQL_READ_PORT"`
	DBMySQLReadUser     string `mapstructure:"DB_MYSQL_READ_USER"`
	DBMySQLReadPassword string `mapstructure:"DB_MYSQL_READ_PASSWORD"`
	DBMySQLReadName     string `mapstructure:"DB_MYSQL_READ_NAME"`
	DBMySQLReadTimezone string `mapstructure:"DB_MYSQL_READ_TIMEZONE"`

	// ======================
	// MYSQL WRITE
	// ======================
	DBMySQLWriteHost     string `mapstructure:"DB_MYSQL_WRITE_HOST"`
	DBMySQLWritePort     string `mapstructure:"DB_MYSQL_WRITE_PORT"`
	DBMySQLWriteUser     string `mapstructure:"DB_MYSQL_WRITE_USER"`
	DBMySQLWritePassword string `mapstructure:"DB_MYSQL_WRITE_PASSWORD"`
	DBMySQLWriteName     string `mapstructure:"DB_MYSQL_WRITE_NAME"`
	DBMySQLWriteTimezone string `mapstructure:"DB_MYSQL_WRITE_TIMEZONE"`

	// ======================
	// POSTGRES
	// ======================
	DBPostgresMaxRetry      int `mapstructure:"DB_POSTGRES_MAX_RETRY"`
	DBPostgresRetryWaitTime int `mapstructure:"DB_POSTGRES_RETRY_WAIT_TIME"`

	DBPostgresReadHost     string `mapstructure:"DB_POSTGRES_READ_HOST"`
	DBPostgresReadPort     string `mapstructure:"DB_POSTGRES_READ_PORT"`
	DBPostgresReadUser     string `mapstructure:"DB_POSTGRES_READ_USER"`
	DBPostgresReadPassword string `mapstructure:"DB_POSTGRES_READ_PASSWORD"`
	DBPostgresReadName     string `mapstructure:"DB_POSTGRES_READ_NAME"`
	DBPostgresReadTimezone string `mapstructure:"DB_POSTGRES_READ_TIMEZONE"`

	DBPostgresWriteHost     string `mapstructure:"DB_POSTGRES_WRITE_HOST"`
	DBPostgresWritePort     string `mapstructure:"DB_POSTGRES_WRITE_PORT"`
	DBPostgresWriteUser     string `mapstructure:"DB_POSTGRES_WRITE_USER"`
	DBPostgresWritePassword string `mapstructure:"DB_POSTGRES_WRITE_PASSWORD"`
	DBPostgresWriteName     string `mapstructure:"DB_POSTGRES_WRITE_NAME"`
	DBPostgresWriteTimezone string `mapstructure:"DB_POSTGRES_WRITE_TIMEZONE"`
}

var (
	conf        Config
	once        sync.Once
	initialized bool
)

func Init(log *log.AppLog) error {
	var err error

	once.Do(func() {

		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		if err = viper.ReadInConfig(); err != nil {
			log.WithError(err).Warn("Could not load .env file, continuing with system environment")
		} else {
			log.Info("Successfully loaded .env file")
		}

		if err = viper.Unmarshal(&conf); err != nil {
			log.WithError(err).Error("Failed to unmarshal configuration")
			return
		}

		// Basic validation
		if conf.AppHost == "" {
			log.Error("APP_HOST is required but not set")
		}

		if conf.AppPort == 0 {
			log.Error("APP_PORT is required but not set")
		}

		initialized = true
		log.Info("Configuration initialized successfully")
	})

	return err
}

func Get(l *log.AppLog) *Config {
	if !initialized {
		if err := Init(l); err != nil {
			l.WithError(err).Fatal("Failed to initialize configuration")
		}
	}
	return &conf
}

func (c *Config) Debug(l *log.AppLog) {
	l.WithField("config", c).Debug("Current configuration")
}
