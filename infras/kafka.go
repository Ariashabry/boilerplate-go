package infras

import (
	"context"
	"fmt"

	applog "github.com/ariashabry/boilerplate-go/helpers/log"

	"github.com/ariashabry/boilerplate-go/helpers/env"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

// RedisNewClient create new instance of redis
func ProvideRedis(config *env.Config, log *applog.AppLog) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.CacheRedisPrimaryHost, config.CacheRedisPrimaryPort),
		Password: config.CacheRedisPrimaryPassword,
		DB:       config.CacheRedisPrimaryDB,
	})

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	log.WithFields(map[string]interface{}{
		"name":   "redis",
		"host":   config.CacheRedisPrimaryHost,
		"port":   config.CacheRedisPrimaryPort,
		"dbName": config.CacheRedisPrimaryDB,
	}).Info("Connected to Redis successfully")

	return &Redis{
		Client: client,
	}
}
