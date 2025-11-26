package main

import (
	"fmt"
	"log"
	"moneybkd/config"
	"moneybkd/controllers"
	"moneybkd/repository"
	"moneybkd/service"
	"net/http"

	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	config.ConnectMongo("mongodb://localhost:27017", "currencydb")

	countryRepo := repository.NewCountryRepository(config.Mongo)
	countryHistory := repository.NewHistoryRepository(config.Mongo)

	apiKey := "fca_live_gd2iQfq2GqTj6sifNV284xdYh6ekp53gXkOhMwjB"
	svc := service.NewCurrencyService(countryRepo, countryHistory, apiKey)
	ctrl := controllers.NewCurrencyController(svc)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"data": "Data",
		})
	})
	e.GET("/currency/:code", ctrl.GetCurrency)
	e.GET("/currency/:code/history", ctrl.GetHistory)
	e.GET("/countries", ctrl.GetCountries)

	log.Println("InitInsert cron...")
	if err := svc.InitCountriesInsert(); err != nil {
		log.Fatalf("Init countries failed %s", err.Error())
	}

	go svc.UpdateDB()

	log.Println("Starting cron...")
	if err := StartCron(svc); err != nil {
		log.Fatalf("failed to start crons: %s", err.Error())
	}

	e.Start(":9992")
}

func StartCron(svc service.CurrencyService) error {
	c, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("cron failed %s", err.Error())
	}

	s := "0 1,13 * * *"
	log.Printf("starting UpdateDB with cron: %s", s)
	job, err := c.NewJob(gocron.CronJob(s, false), gocron.NewTask(svc.UpdateDB))

	if err != nil {
		return fmt.Errorf("failed to create UpdateDB job %s", err.Error())
	}

	log.Printf("Starting crons %s %s", job.Name(), job.ID())
	c.Start()
	return nil
}
