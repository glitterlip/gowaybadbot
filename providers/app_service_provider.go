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
		Format:           `[${time_custom}] ["${remote_ip}"] [${status}] [${method}] [${uri}] [${error}] [time:${latency_human}] [in:${bytes_in}] [out:${bytes_out}]` + "\n" + `[ua:${user_agent}]` + "\n",
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
