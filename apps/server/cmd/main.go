package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antoniosarro/gosvelte/apps/server/config"
	"github.com/antoniosarro/gosvelte/apps/server/internal/server"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/db"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/logger"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/redis"

	"github.com/antoniosarro/gosvelte/apps/server/docs"
)

//	@title			gosvelte monorepo
//	@version		0.0.1
//	@description	Boilerplate for Echo Golang development.

//	@contact.name	Antonio Sarro
//	@contact.url	https://www.antoniosarro.dev
//	@contact.email	contact@antoniosarro.dev

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @BasePath	/api/v1
func main() {
	conf, err := config.Load()
	if err != nil {
		panic("load config error, " + err.Error())
	}

	log := logger.Init(conf)
	ctx := context.Background()

	if err := run(ctx, conf, log); err != nil {
		log.Fatalf("run error, %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, conf *config.Config, log *logger.Log) error {
	log.Infof("starting server...")

	db, err := db.Init(conf)
	if err != nil {
		log.Fatalf("database connection error, %v", err)
	}
	defer db.Close()
	log.Infof("database connected: %+v", db.Stats())

	rdb, err := redis.Init(conf)
	if err != nil {
		log.Fatalf("redis connection error, %v", err)
	}
	defer rdb.Close()
	log.Infof("redis connected: %+v", rdb.PoolStats())

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	server := server.Init(&server.Options{
		DB:       db,
		Cache:    rdb,
		Log:      log,
		ServConf: conf,
	})

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
	serverErrCh := make(chan error, 1)

	go func() {
		s := &http.Server{
			Addr:         conf.Server.Host + ":" + conf.Server.Port,
			ReadTimeout:  time.Second * conf.Server.ReadTimeout,
			WriteTimeout: time.Second * conf.Server.WriteTimeout,
		}

		log.Infof("server started on %s:%s", conf.Server.Host, conf.Server.Port)
		serverErrCh <- server.StartServer(s)
	}()

	select {
	case sig := <-shutdownCh:
		log.Infof("shutdown started: %s", sig)
		defer log.Infof("shutdown completed: %s", sig)

		ctx, cancel := context.WithTimeout(ctx, conf.Server.CtxDefaultTimeout*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("gracefull shutdown failed, server forced to shutdown: %v", err)
		}
	case err := <-serverErrCh:
		log.Errorf("server error, %v", err)
		return err
	}

	return nil
}
