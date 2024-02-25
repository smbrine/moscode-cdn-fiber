package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron/v3"
	"log"
	"moscode-cdn-fiber/api/handlers"
	"moscode-cdn-fiber/configs"
	cdnCron "moscode-cdn-fiber/internal/cron"
)

func main() {
	app := fiber.New()
	c := cron.New(cron.WithSeconds())

	appConfig := configs.GetConfig()

	staticSpec := fmt.Sprintf("*/%v * * * * *", appConfig.RefreshFilesSec)
	indexSpec := fmt.Sprintf("*/%v * * * * *", appConfig.RefreshIndexSec)

	cdnCron.ScheduleCronJob(c, staticSpec, cdnCron.UpdateStaticJob)
	cdnCron.ScheduleCronJob(c, indexSpec, cdnCron.UpdateIndexJob)

	app.Use(cors.New())

	for _, page := range configs.UrlPages {
		app.Get(page, func(c *fiber.Ctx) error { c.Path("/index.html"); return handlers.HandleStatic(c) })
	}
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(&fiber.Map{
			"health": "ok",
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return handlers.HandleStatic(c)
	})

	c.Start()
	defer c.Stop()
	addr := appConfig.ServerHost
	port := appConfig.ServerPort

	log.Fatal(app.Listen(fmt.Sprintf("%v:%v", addr, port)))
}
