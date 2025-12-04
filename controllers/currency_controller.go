package controllers

import (
	"log"
	"moneybkd/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CurrencyController struct {
	svc service.CurrencyService
}

func NewCurrencyController(s service.CurrencyService) *CurrencyController {
	return &CurrencyController{svc: s}
}

func (ctr *CurrencyController) GetCurrency(c echo.Context) error {
	code := c.Param("code")

	res, err := ctr.svc.GetCurrency(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr *CurrencyController) GetCountries(c echo.Context) error {
	res, err := ctr.svc.GetCountries(c.Request().Context())

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr *CurrencyController) GetHistory(ctx echo.Context) error {
	code := ctx.Param("code")

	log.Println("GET history 1")
	h, err := ctr.svc.GetHistoryByCode(ctx.Request().Context(), code)
	log.Println("GET history 1")
	log.Println(h)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, h)
}
