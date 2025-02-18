package db

import (
	"fmt"
	"time"

	"github.com/antoniosarro/gosvelte/apps/server/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func Init(cfg *config.Config) (*sqlx.DB, error) {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.DbName,
		cfg.DB.SSLMode,
	)

	db, err := sqlx.Connect(cfg.DB.Driver, connection)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIddleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.DB.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(cfg.DB.ConnMaxIddleTime) * time.Second)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
