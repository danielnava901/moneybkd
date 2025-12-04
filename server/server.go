package server

import (
	"moneybkd/config"
	"moneybkd/controllers"
	"moneybkd/repository"
	"moneybkd/service"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
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

	return e
}
