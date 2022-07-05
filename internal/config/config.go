// Package define types for application configuration.
// Also has method for configuration parsing and loading.
// Uses github.com/cristalhq/aconfig package for configuration parsing.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

// Env defines application environment type.
type Env string

const (
	// EnvLocal defines local environment.
	EnvLocal Env = "local"
	// EnvProduction defines production environment.
	EnvProduction Env = "production"
	// EnvStaging defines staging environment.
	EnvStaging Env = "staging"
)

// Config define application configuration.
type Config struct {
	Env      Env `default:"local" usage:"app environment (local, production, staging)"`
	Log      Log
	Sentry   Sentry
	Postgres Postgres
	Bot      Bot
	HTTP     HTTP
}

// Load application configuration from following sources:
// - environment variables with prefix "GO_TG_BOT_*"
// - command line flags
// - yaml files
func Load(files []string) *Config {
	cfg := Config{}

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		EnvPrefix: "GO_TG_BOT",
		Files:     files,
		FileDecoders: map[string]aconfig.FileDecoder{
			".yml": aconfigyaml.New(),
		},
	})

	if err := loader.Load(); errors.Is(err, flag.ErrHelp) {
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(2)
	}

	return &cfg
}
