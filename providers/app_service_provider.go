package providers

import (
	"github.com/coocood/freecache"
	"github.com/labstack/echo/v4"
	"goawaybot/store"
	"runtime/debug"
)

func RegisterAppServiceProvider(e *echo.Echo) {
	e.IPExtractor = echo.ExtractIPDirect()

	RegisterCacheService()
}
func RegisterCacheService() {
	cacheSize := 100 * 1024 * 1024
	store.Cache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)

}
