package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"moscode-cdn-fiber/configs"
	"sync"
)

var (
	redisConnect     *redis.Options
	redisConnectOnce sync.Once
)

func getRedisOpt() (*redis.Options, error) {
	var err error
	redisConnectOnce.Do(func() {
		appConfig := configs.GetConfig()
		if appConfig.RedisURL != "" {

			if a, err := redis.ParseURL(appConfig.RedisURL); err == nil {
				redisConnect = a
			} else {
				log.Fatal(err)
			}
		} else {
			err = errors.New("redis url can't be empty")
		}
	})

	return redisConnect, err
}

func GetRedisClient() (*redis.Client, error) {
	client, err := getRedisOpt()
	return redis.NewClient(client), err
}
