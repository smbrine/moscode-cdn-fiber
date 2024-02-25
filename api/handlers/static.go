package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"log"
	"moscode-cdn-fiber/configs"
	cdnRedis "moscode-cdn-fiber/redis"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func serveFile(
	c *fiber.Ctx,
	client *redis.Client,
	ctx context.Context,
	staticDir, safePath, cacheKey string,
	cacheTime time.Duration,
) error {
	cachedFile, err := client.Get(ctx, cacheKey).Bytes()

	if err == nil && cachedFile != nil {
		return c.Send(cachedFile)
	}

	filePath := filepath.Join(staticDir, safePath)
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}

	client.Set(ctx, cacheKey, fileBytes, cacheTime)

	return c.Send(fileBytes)
}

func serveGzip(
	c *fiber.Ctx,
	client *redis.Client,
	ctx context.Context,
	staticDir, safePath string,
	cacheTime time.Duration,
) error {
	gzipSafePath := safePath + ".gz"
	cacheKey := "file:" + gzipSafePath
	c.Set("Content-Encoding", "gzip")
	return serveFile(c, client, ctx, staticDir, gzipSafePath, cacheKey, cacheTime)
}

func serveBrotli(
	c *fiber.Ctx,
	client *redis.Client,
	ctx context.Context,
	staticDir, safePath string,
	cacheTime time.Duration,
) error {
	brSafePath := safePath + ".br"
	cacheKey := "file:" + brSafePath
	c.Set("Content-Encoding", "br")
	return serveFile(c, client, ctx, staticDir, brSafePath, cacheKey, cacheTime)
}

func serveUncompressed(
	c *fiber.Ctx,
	client *redis.Client,
	ctx context.Context,
	staticDir, safePath string,
	cacheTime time.Duration,
) error {
	cacheKey := "file:" + safePath
	return serveFile(c, client, ctx, staticDir, safePath, cacheKey, cacheTime)
}

func HandleStatic(c *fiber.Ctx) error {
	reqPath := c.Path()
	client, err := cdnRedis.GetRedisClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	appConfig := configs.GetConfig()
	staticDir := appConfig.StaticDir
	cacheFilesFor := time.Duration(appConfig.CacheFilesFor) * time.Second

	safePath := filepath.Clean(reqPath)

	ext := strings.TrimPrefix(filepath.Ext(safePath), ".")
	c.Type(ext)

	acceptEncodingHeader := c.Get("Accept-Encoding")
	if strings.Contains(acceptEncodingHeader, "br") {
		if err := serveBrotli(c, client, ctx, staticDir, safePath, cacheFilesFor); err == nil {
			return nil
		}
	} else if strings.Contains(acceptEncodingHeader, "gzip") {
		if err := serveGzip(c, client, ctx, staticDir, safePath, cacheFilesFor); err == nil {
			return nil
		}
	}

	if err := serveUncompressed(c, client, ctx, staticDir, safePath, cacheFilesFor); err != nil {
		log.Printf("Error serving uncompressed file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return err
}
