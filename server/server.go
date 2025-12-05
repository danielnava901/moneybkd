package server

import (
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

	config.ConnectSupabase()
	countryRepo := repository.NewCountryRepository(config.Supabase)
	countryHistory := repository.NewHistoryRepository(config.Supabase)
	apiKey := os.Getenv("EXCHANGE_API_KEY")

	svc := service.NewCurrencyService(countryRepo, countryHistory, apiKey)
	ctrl := controllers.NewCurrencyController(svc)

	e := echo.New()
	e.Use(middleware.CORS())

	api := e.Group("/api")
	api.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"data": "Data Ok",
		})
	})

	api.GET("/currency/:code", ctrl.GetCurrency)
	api.GET("/currency/:code/history", ctrl.GetHistory)
	api.GET("/countries", ctrl.GetCountries)

	return e
}
