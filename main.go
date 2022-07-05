package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/mr-linch/go-tg-bot/internal/app"
	"github.com/mr-linch/go-tg-bot/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	buildVersion = "unknown"
	buildRef     = "unknown"
	buildTime    = "unknown"
)

func main() {
	// parse cfg
	cfg := config.Load([]string{
		"./go-tg-bot.local.yml",
		"./go-tg-bot.yml",
		"/etc/go-tg-bot.yml",
	})

	// init context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// listen for signals
	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	// add logger to context
	ctx = withLogger(ctx, cfg)

	log.Ctx(ctx).Info().
		Dict("build", zerolog.Dict().
			Str("version", buildVersion).
			Str("ref", buildRef).
			Str("time", buildTime),
		).
		Str("env", string(cfg.Env)).
		Msg("start go-tg-bot")

	if err := app.Run(ctx, cfg, app.BuildInfo{
		Version: buildVersion,
		Ref:     buildRef,
		Time:    buildTime,
	}); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("run failed")
		defer os.Exit(2)
	}
}

func withLogger(ctx context.Context, config *config.Config) context.Context {
	if config.Log.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if config.Log.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return log.Logger.WithContext(ctx)
}
