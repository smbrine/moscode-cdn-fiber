package cron

import (
	"github.com/robfig/cron/v3"
	"log"
)

func ScheduleCronJob(c *cron.Cron, spec string, job func()) {
	job()
	if _, err := c.AddFunc(spec, job); err != nil {
		log.Println("Error scheduling the cron job:", err)
	}
}
