//go:build heroku
// +build heroku

package config

import "os"

func aliasEnv() {
	if os.Getenv("GO_TG_BOT_POSTGRES_DSN") == "" && os.Getenv("POSTGRES_URL") != "" {
		os.Setenv("GO_TG_BOT_POSTGRES_DSN", os.Getenv("POSTGRES_URL"))
	}
}
