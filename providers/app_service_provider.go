package providers

import (
	"github.com/coocood/freecache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"goawaybot/services"
	"goawaybot/store"
	"runtime/debug"
)

func RegisterAppServiceProvider(e *echo.Echo) {
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_custom}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "01-02 15:04:05",
	}))
	RegisterCacheService()
}
func RegisterCacheService() {
	cacheSize := 100 * 1024 * 1024
	store.Cache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
	services.InitStatistics()
}
