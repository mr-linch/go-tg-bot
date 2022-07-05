package store

import (
	"context"

	"github.com/friendsofgo/errors"
)

var (
	// ErrTooManyAffectedRows returned when store modifes too many rows.
	ErrTooManyAffectedRows = errors.New("too many affected rows")
)

// MigratorFactory define method for migrations of store
type MigratorFactory interface {
	Migrator() Migrator
}

// Migrator defines generic interface for migrations.
type Migrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}

// StoreFactory define interface of factory methods
type StoreFactory interface {
	User() User
}

// StoreTx define interface of transactional of store.
type StoreTx interface {
	// TxFactory returns function for create transaction scopes.
	Tx(ctx context.Context, txFunc TxFunc) error
}

// TxFunc define signature of callback used in tx block
type TxFunc func(ctx context.Context) error

// Txier define function to start tx block
type Txier func(ctx context.Context, txFunc TxFunc) error

//go:generate mockery --name Store  --case underscore

// Store define generic interface for database with transaction support
type Store interface {
	StoreFactory
	StoreTx
	MigratorFactory
}

func Asc(field string) string {
	return field + " ASC"
}

func Desc(field string) string {
	return field + " DESC"
}
