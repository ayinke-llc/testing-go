package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	LogLevel string `json:"level,omitempty" env:"LOG_LEVEL"`

	PostgresDSN string `env:"POSTGRES_DSN" json:"postgres_dsn,omitempty"`

	// DISABLE IN PROD
	PostgresLogQueries bool `env:"POSTGRES_LOG_QUERIES" json:"postgres_log_queries,omitempty"`

	RedisDSN string `env:"REDIS_DSN" json:"redis_dsn,omitempty"`
}

func Load() (Config, error) {
	var cfg Config

	return cfg, env.Parse(&cfg)
}
