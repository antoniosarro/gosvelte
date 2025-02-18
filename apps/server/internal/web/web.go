package web

import (
	"net/http"

	"github.com/antoniosarro/gosvelte/apps/server/config"
	"github.com/antoniosarro/gosvelte/apps/server/internal/web/middlewares"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

type Web struct {
	Echo *echo.Echo
	log  *logger.Log
	Mid  *middlewares.Middleware
}

func New(log *logger.Log) *Web {
	return &Web{
		Echo: echo.New(),
		log:  log,
	}
}

func (w *Web) InitCustomMiddleware(servConf *config.Config, cache *redis.Client) {
	w.Mid = middlewares.New(servConf, w.log, cache)
}

func (w *Web) EnableCORSMware(origins []string) {
	w.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (w *Web) EnableRecovererMware() {
	w.Echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
}

func (w *Web) EnableGlobalMiddleware() {
	w.Echo.Use(middleware.RequestID())
	w.Echo.Use(middleware.Secure())
	w.Echo.Use(middleware.BodyLimit("10M"))
	w.Echo.Use(w.Mid.RequestLoggerMware)
	w.Echo.Use(w.Mid.ErrorLoggingMware)
}
