package controllers

import (
	"github.com/labstack/echo/v4"
	"goawaybot/services"
	"net/http"
)

func Index(c echo.Context) error {
	m := services.GetSystemStatistics()
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"m": m,
	})
}
