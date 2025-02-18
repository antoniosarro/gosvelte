package server

import (
	"github.com/antoniosarro/gosvelte/apps/server/config"
	"github.com/antoniosarro/gosvelte/apps/server/internal/web"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	DB       *sqlx.DB
	Cache    *redis.Client
	Log      *logger.Log
	ServConf *config.Config
}

func Init(opts *Options) *echo.Echo {
	w := web.New(opts.Log)

	w.Echo.HideBanner = true
	w.Echo.HidePort = true

	// middleware setup
	w.InitCustomMiddleware(opts.ServConf, opts.Cache)
	w.EnableCORSMware(opts.ServConf.Server.AllowedOrigins)
	w.EnableRecovererMware()
	w.EnableGlobalMiddleware()

	// remap all routes
	router(w, opts)

	return w.Echo
}
