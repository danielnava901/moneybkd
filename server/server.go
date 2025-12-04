package server

import (
	"log"
	"moneybkd/config"
	"moneybkd/controllers"
	"moneybkd/repository"
	"moneybkd/service"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	log.Println("INIT NEW SERVER..............1.")
	config.ConnectSupabase()
	countryRepo := repository.NewCountryRepository(config.Supabase)
	countryHistory := repository.NewHistoryRepository(config.Supabase)
	apiKey := os.Getenv("EXCHANGE_API_KEY")
	log.Println("INIT NEW SERVER..............2.")
	svc := service.NewCurrencyService(countryRepo, countryHistory, apiKey)
	ctrl := controllers.NewCurrencyController(svc)

	log.Println("INIT NEW SERVER..............3.")
	e := echo.New()
	e.Use(middleware.CORS())

	log.Println("INIT NEW SERVER..............4.")
	api := e.Group("/api")

	api.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"data": "Data Ok",
		})
	})

	log.Println("INIT NEW SERVER..............5.")
	api.GET("/currency/:code", ctrl.GetCurrency)
	api.GET("/currency/:code/history", ctrl.GetHistory)
	api.GET("/countries", ctrl.GetCountries)

	return e
}
