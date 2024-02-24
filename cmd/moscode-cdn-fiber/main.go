package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron/v3"
	iLog "log"
	"moscode-cdn-fiber/api/handlers"
	"moscode-cdn-fiber/configs"
	cron2 "moscode-cdn-fiber/internal/cron"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	cron2.UpdateStaticJob()
	cron2.UpdateIndexJob()

	for _, page := range configs.UrlPages {
		app.Get(page, func(c *fiber.Ctx) error { c.Path("/index.html"); return handlers.HandleStatic(c) })
	}

	app.Use(func(c *fiber.Ctx) error {
		return handlers.HandleStatic(c)
	})

	c := cron.New(cron.WithSeconds())

	if _, err := c.AddFunc("*/10 * * * * *", cron2.UpdateStaticJob); err != nil {
		iLog.Println("Error scheduling the cron job:", err)
	}

	if _, err := c.AddFunc("*/10 * * * * *", cron2.UpdateIndexJob); err != nil {
		iLog.Println("Error scheduling the cron job:", err)
	}

	c.Start()
	defer c.Stop()

	log.Fatal(app.Listen(":8080"))
}
