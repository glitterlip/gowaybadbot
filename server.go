package main

import (
	"embed"
	"github.com/labstack/echo/v4"
	"goawaybot/http/controllers"
	"goawaybot/providers"
	"goawaybot/services"
)

//go:embed  resources images
var appFs embed.FS

func main() {
	providers.TemplateFs = appFs
	services.ImagesFs = appFs

	e := echo.New()
	providers.RegisterAppServiceProvider(e)
	providers.RegisterViewServiceProvider(e)
	//static
	e.StaticFS("/static", echo.MustSubFS(appFs, "resources/static"))
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "^-^")
	})
	e.GET("/challenge", controllers.Challenge).Name = "Challenge"
	e.POST("/check", controllers.Check).Name = "check"
	e.GET("/verify", controllers.Verify).Name = "verify"
	e.GET("/refresh", controllers.Refresh).Name = "refresh"
	e.Logger.Fatal(e.Start(":11323"))

}
