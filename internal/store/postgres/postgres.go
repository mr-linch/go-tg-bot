package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/friendsofgo/errors"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/mr-linch/go-tg-bot/internal/store/postgres/migrations"
	"github.com/mr-linch/go-tg-bot/internal/store/postgres/shared"
	"github.com/rs/zerolog/log"
)

//go:generate sqlboiler psql

type Postgres struct {
	*sql.DB
	migrator *migrations.Migrator

	user *userStore
}

var _ store.Store = &Postgres{}

// New create postgres based database with all stores.
func New(db *sql.DB) *Postgres {
	pg := &Postgres{
		DB:       db,
		migrator: migrations.New(db),
	}

	base := baseStore{DB: db, Txier: pg.Tx}

	pg.user = &userStore{base}

	return pg
}

func (pg *Postgres) User() store.User {
	return pg.user
}

func (pg *Postgres) Migrator() store.Migrator {
	return pg.migrator
}

// Tx run code in database transaction.
// Based on: https://stackoverflow.com/a/23502629.
func (pg *Postgres) Tx(ctx context.Context, txFunc store.TxFunc) (err error) {
	tx := shared.GetTx(ctx)

	if tx != nil {
		return txFunc(ctx)
	}

	log.Ctx(ctx).Trace().Msg("start transaction")

	tx, err = pg.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin tx failed")
	}

	ctx = shared.WithTx(ctx, tx)

	//nolint:gocritic
	defer func() {
		if r := recover(); r != nil {
			log.Ctx(ctx).Trace().Err(err).Msg("rollback transaction")
			if err := tx.Rollback(); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("transaction rollback failed")
			}
			panic(r)
		} else if err != nil {
			log.Ctx(ctx).Trace().Err(err).Msg("rollback transaction")
			if err := tx.Rollback(); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("transaction rollback failed")
			}
		} else {
			log.Ctx(ctx).Trace().Msg("commit transaction")
			err = tx.Commit()
			if err != nil {
				log.Ctx(ctx).Trace().Err(err).Msg("commit failed")
			}
		}
	}()

	err = txFunc(ctx)

	return err
}
