package app

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/getsentry/sentry-go"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/config"
	"github.com/mr-linch/go-tg-bot/internal/delivery/bot"
	"github.com/mr-linch/go-tg-bot/internal/service/container"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/mr-linch/go-tg-bot/internal/store/postgres"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/rs/zerolog/log"
)

type BuildInfo struct {
	Version string
	Ref     string
	Time    string
}

func Run(ctx context.Context, cfg *config.Config, buildInfo BuildInfo) error {
	close, err := initSentry(ctx, cfg, buildInfo)
	if err != nil {
		return errors.Wrap(err, "init sentry")
	}
	defer close()

	store, close, err := newStore(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "init store")
	}
	defer close()

	botClient, err := newBotClient(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "init bot client")
	}

	srv := container.New(container.Deps{
		Store: store,
	})

	bot := bot.New(&bot.Deps{
		Service: srv,
	})

	if cfg.Bot.Webhook.BaseURL != "" {
		fullURL := strings.TrimRight(cfg.Bot.Webhook.BaseURL, "/") + cfg.Bot.Webhook.Path

		log.Ctx(ctx).Info().
			Str("url", fullURL).
			Str("listen", cfg.HTTP.Listen).
			Msg("start webhook server...")

		return tgb.NewWebhook(
			bot,
			botClient,
			fullURL,
		).Run(
			ctx,
			cfg.HTTP.Listen,
		)
	} else {
		log.Ctx(ctx).Info().Msg("start polling...")

		return tgb.NewPoller(
			bot,
			botClient,
		).Run(
			ctx,
		)
	}
}

func initSentry(ctx context.Context, cfg *config.Config, build BuildInfo) (context.CancelFunc, error) {
	if cfg.Sentry.DSN != "" {
		log.Ctx(ctx).Info().Str("dsn", cfg.Sentry.DSN).Msg("init sentry...")
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry.DSN,
			Environment:      string(cfg.Env),
			AttachStacktrace: true,
			Release:          build.Version + "-" + build.Ref,
			// Debug:            true,
		}); err != nil {
			return nil, err
		}

		defer sentry.Flush(2 * time.Second)

		return func() {
			log.Ctx(ctx).Debug().Msg("flush sentry events")
			sentry.Flush(2 * time.Second)
		}, nil
	} else {
		log.Ctx(ctx).Warn().Msg("sentry is disabled...")
	}

	return func() {}, nil
}

func newStore(ctx context.Context, cfg *config.Config) (store.Store, context.CancelFunc, error) {
	// open and ping db
	log.Ctx(ctx).Info().
		Str("dsn", cfg.Postgres.DSN).
		Int("max_open_conns", cfg.Postgres.MaxOpenConns).
		Int("max_idle_conns", cfg.Postgres.MaxIdleConns).
		Msg("open db")

	db, err := sql.Open("postgres", cfg.Postgres.DSN)
	if err != nil {
		return nil, nil, errors.Wrap(err, "open postgres")
	}

	close := func() {
		db.Close()
	}

	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)

	log.Ctx(ctx).Debug().Msg("ping db...")
	started := time.Now()
	if err := db.PingContext(ctx); err != nil {
		close()
		return nil, nil, errors.Wrap(err, "ping db")
	}
	log.Ctx(ctx).Debug().Dur("took", time.Since(started)).Msg("ping db - done")

	pg := postgres.New(db)

	if err := pg.Migrator().Up(ctx); err != nil {
		close()
		return nil, nil, errors.Wrap(err, "migrate db")
	}

	return pg, func() {
		log.Ctx(ctx).Debug().Msg("close db...")
		db.Close()
	}, nil
}

func newBotClient(ctx context.Context, config *config.Config) (*tg.Client, error) {
	client := tg.New(config.Bot.Token)

	me, err := client.Me(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get me")
	}

	log.Printf("authorized as bot https://t.me/%s", me.Username)

	return client, nil
}
