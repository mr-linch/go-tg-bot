package postgres

import (
	"context"
	"testing"

	"github.com/mr-linch/go-tg-bot/internal/store/postgres/postgrestest"
)

func newPostgres(t *testing.T) (context.Context, *Postgres) {
	t.Helper()

	ctx := context.Background()

	db := postgrestest.New(t)

	pg := New(db)

	if err := pg.Migrator().Up(ctx); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	return ctx, pg
}
