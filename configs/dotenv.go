package configs

import (
	"log"
	"os"
	"strconv"
	"sync"
)

type AppConfig struct {
	TempDir         string
	StaticDir       string
	RedisURL        string
	BaseURL         string
	ServerHost      string
	ServerPort      int
	RefreshFilesSec int
	RefreshIndexSec int
	CacheFilesFor   int
}

var config *AppConfig
var once sync.Once

func GetConfig() *AppConfig {
	once.Do(func() {
		config = &AppConfig{
			TempDir:         os.Getenv("TEMP_DIR"),
			StaticDir:       os.Getenv("STATIC_DIR"),
			RedisURL:        os.Getenv("REDIS_CDN_URL"),
			BaseURL:         os.Getenv("JS_SSR_URL"),
			ServerHost:      os.Getenv("SERVER_HOST"),
			ServerPort:      atoiOrFallback(os.Getenv("SERVER_PORT"), 8080),
			RefreshFilesSec: atoiOrFallback(os.Getenv("REFRESH_FILES_SEC"), 30),
			RefreshIndexSec: atoiOrFallback(os.Getenv("REFRESH_INDEX_SEC"), 15),
			CacheFilesFor:   atoiOrFallback(os.Getenv("CACHE_FILES_FOR"), 15),
		}
	})
	return config
}

func atoiOrFallback(value string, fallback int) int {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	} else {
		log.Printf("Error converting %s to int: %v", value, err)
		return fallback
	}
}
