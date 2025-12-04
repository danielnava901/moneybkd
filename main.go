package main

import (
	"fmt"
	"log"
	"moneybkd/server"
	"moneybkd/service"

	"github.com/go-co-op/gocron/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var echoServer = server.New()
	log.Println("Starting server in :9992")
	echoServer.Start(":9992")
}

func StartCron(svc service.CurrencyService) error {
	c, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("cron failed %s", err.Error())
	}

	s := "0 1,7,13,19 * * *"
	log.Printf("starting UpdateDB with cron: %s", s)

	job, err := c.NewJob(gocron.CronJob(s, false), gocron.NewTask(svc.UpdateDB))

	if err != nil {
		return fmt.Errorf("failed to create UpdateDB job %s", err.Error())
	}

	log.Printf("Starting crons %s %s", job.Name(), job.ID())
	c.Start()
	return nil
}
