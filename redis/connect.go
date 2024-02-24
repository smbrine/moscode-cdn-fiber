package redis

import (
	redisDB "github.com/redis/go-redis/v9"
	"log"
	"moscode-cdn-fiber/configs"
	"sync"
)

var (
	redisConnect     *redisDB.Options
	redisConnectOnce sync.Once
)

func getRedisOpt() *redisDB.Options {
	redisConnectOnce.Do(func() {
		appConfig := configs.GetConfig()
		if a, e := redisDB.ParseURL(appConfig.RedisURL); e == nil {
			redisConnect = a
		} else {
			log.Fatal(e)
		}
	})

	return redisConnect
}

func GetRedisClient() *redisDB.Client {
	return redisDB.NewClient(getRedisOpt())
}
