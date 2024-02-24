package configs

import (
	"os"
	"sync"
)

type DotEnv struct {
	TempDir   string
	StaticDir string
	RedisURL  string
}

var (
	TempDir       string
	StaticDir     string
	RedisURL      string
	BaseURL       string
	TempDirOnce   sync.Once
	StaticDirOnce sync.Once
	RedisURLOnce  sync.Once
	BaseURLOnce   sync.Once
)

func GetTempDir() string {
	TempDirOnce.Do(func() {
		TempDir = os.Getenv("TEMP_DIR")
	})
	return TempDir
}

func GetStaticDir() string {
	StaticDirOnce.Do(func() {
		StaticDir = os.Getenv("STATIC_DIR")
	})
	return StaticDir
}

func GetRedisURL() string {
	RedisURLOnce.Do(func() {
		RedisURL = os.Getenv("REDIS_CDN_URL")
	})
	return RedisURL
}

func GetBaseURL() string {
	BaseURLOnce.Do(func() {
		BaseURL = os.Getenv("JS_SSR_URL")
	})
	return BaseURL
}
