//go:build heroku
// +build heroku

package config

import "os"

func setRuntime(cfg *Config) {
	cfg.Postgres.DSN = os.Getenv("DATABASE_URL")
}
