package redis

import (
	"github.com/antoniosarro/gosvelte/apps/server/config"
	"github.com/redis/go-redis/v9"
)

func Init(conf *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Cache.Url,
		Password: "",
		DB:       0,
		Protocol: 3,
	})

	return rdb, nil
}
