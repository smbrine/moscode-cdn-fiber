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

func serveGzipFile(
	c *fiber.Ctx,
	client *redis.Client,
	ctx context.Context,
	staticDir, safePath string,
	cacheTime time.Duration,
) error {
	gzipSafePath := safePath + ".gz"
	cacheKey := "file:" + gzipSafePath

	cachedFile, err := client.Get(ctx, cacheKey).Bytes()
	if err == nil && cachedFile != nil {
		c.Set("Content-Encoding", "gzip")
		ext := strings.TrimPrefix(filepath.Ext(safePath), ".")
		c.Type(ext)
		return c.Send(cachedFile)
	}

	compressedFilePath := filepath.Join(staticDir, gzipSafePath)
	fileBytes, err := os.ReadFile(compressedFilePath)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Compressed file not found")
	}

	client.Set(ctx, cacheKey, fileBytes, cacheTime)

	c.Set("Content-Encoding", "gzip")
	ext := strings.TrimPrefix(filepath.Ext(safePath), ".")
	c.Type(ext)
	return c.Send(fileBytes)
}

func serveBrotliFile(
	c *fiber.Ctx,
	client *redis.Client,
	ctx context.Context,
	staticDir, safePath string,
	cacheTime time.Duration,
) error {
	brotliSafePath := safePath + ".br"
	cacheKey := "file:" + brotliSafePath

	cachedFile, err := client.Get(ctx, cacheKey).Bytes()
	if err == nil && cachedFile != nil {
		c.Set("Content-Encoding", "br")
		ext := strings.TrimPrefix(filepath.Ext(safePath), ".")
		c.Type(ext)
		return c.Send(cachedFile)
	}

	compressedFilePath := filepath.Join(staticDir, brotliSafePath)
	fileBytes, err := os.ReadFile(compressedFilePath)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Compressed file not found")
	}

	client.Set(ctx, cacheKey, fileBytes, cacheTime)

	c.Set("Content-Encoding", "br")
	ext := strings.TrimPrefix(filepath.Ext(safePath), ".")
	c.Type(ext)
	return c.Send(fileBytes)
}

func serveUncompressedFile(
	c *fiber.Ctx,
	client *redis.Client,
	ctx context.Context,
	staticDir, safePath string,
	cacheTime time.Duration,
) error {
	cacheKey := "file:" + safePath

	cachedFile, err := client.Get(ctx, cacheKey).Bytes()
	if err == nil && cachedFile != nil {
		ext := strings.TrimPrefix(filepath.Ext(safePath), ".")
		c.Type(ext)
		return c.Send(cachedFile)
	}

	filePath := filepath.Join(staticDir, safePath)
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}

	client.Set(ctx, cacheKey, fileBytes, cacheTime)

	ext := strings.TrimPrefix(filepath.Ext(safePath), ".")
	c.Type(ext)

	return c.Send(fileBytes)
}

func HandleStatic(c *fiber.Ctx) error {
	reqPath := c.Path()
	client := cdnRedis.GetRedisClient()
	ctx := context.Background()
	staticDir := configs.GetStaticDir()

	safePath := filepath.Clean(reqPath)
	acceptEncodingHeader := c.Get("Accept-Encoding")

	if strings.Contains(acceptEncodingHeader, "br") {
		if err := serveBrotliFile(c, client, ctx, staticDir, safePath, 10*time.Second); err == nil {
			return nil
		}
	} else if strings.Contains(acceptEncodingHeader, "gzip") {
		if err := serveGzipFile(c, client, ctx, staticDir, safePath, 10*time.Second); err == nil {
			return nil
		}
	}

	if err := serveUncompressedFile(c, client, ctx, staticDir, safePath, 10*time.Second); err != nil {
		log.Printf("Error serving uncompressed file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
