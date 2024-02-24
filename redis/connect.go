package redis

import (
	redisDB "github.com/redis/go-redis/v9"
	cdnConfig "moscode-cdn-fiber/configs"
	"sync"
)

var (
	redisConnect     *redisDB.Options
	redisConnectOnce sync.Once
)

func getRedisOpt() *redisDB.Options {
	redisConnectOnce.Do(func() {
		if a, e := redisDB.ParseURL(cdnConfig.GetRedisURL()); e == nil {
			redisConnect = a
		} else {
			println(e)
		}
	})

	return redisConnect
}

func GetRedisClient() *redisDB.Client {
	return redisDB.NewClient(getRedisOpt())
}
