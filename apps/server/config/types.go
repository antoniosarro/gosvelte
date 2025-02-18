package config

import "time"

type Config struct {
	App    string `env:"APP_ENV"`
	Server serverConfig
	DB     postgresConfig
	Cache  redisConfig
	Logger loggerConfig
}

type serverConfig struct {
	Host              string        `env:"SERVER_HOST"`
	Port              string        `env:"SERVER_PORT"`
	ReadTimeout       time.Duration `env:"SERVER_READTIMEOUT"`
	WriteTimeout      time.Duration `env:"SERVER_WRITETIMEOUT"`
	CtxDefaultTimeout time.Duration `env:"SERVER_CTXDEFAULTTIMEOUT"`
	AllowedOrigins    []string      `env:"SERVER_ALLOWEDORIGINS"`
	JWTSecret         string        `env:"SERVER_JWTSECRET"`
}

type loggerConfig struct {
	Level          string `env:"LOGGER_LEVEL"`
	File           string `env:"LOGGER_FILE"`
	FileMaxSize    int    `env:"LOGGER_FILE_MAXSIZE"`
	FileMaxBackups int    `env:"LOGGER_FILE_MAXBACKUPS"`
	FileMaxAge     int    `env:"LOGGER_FILE_MAXAGE"`
	FileCompress   bool   `env:"LOGGER_FILE_COMPRESS"`
}

type postgresConfig struct {
	Host             string `env:"DB_HOST"`
	Port             string `env:"DB_PORT"`
	User             string `env:"DB_USER"`
	Password         string `env:"DB_PASSWORD"`
	DbName           string `env:"DB_NAME"`
	SSLMode          string `env:"DB_SSLMODE"`
	Driver           string `env:"DB_DRIVER"`
	MaxIddleConns    int    `env:"DB_MAXIDLECONNS"`
	MaxOpenConns     int    `env:"DB_MAXOPENCONNS"`
	ConnMaxIddleTime int8   `env:"DB_CONNMAXIDLETIME"`
	ConnMaxLifetime  int8   `env:"DB_CONNMAXLIFETIME"`
}

type redisConfig struct {
	Url string `env:"REDIS_URL"`
}
