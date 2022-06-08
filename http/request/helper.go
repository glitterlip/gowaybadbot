package request

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func IsJosn(c *echo.Context) bool {
	want := (*c).Request().Header.Get("Accept")
	if strings.Contains(strings.ToLower(want), "json") {
		return true
	}
	return false
}
