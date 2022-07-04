package store

import "context"

type Base[M any, Q any] interface {
	// Add T to store.
	Add(ctx context.Context, v *M) error

	// Update T in store.
	Update(ctx context.Context, v *M, fields ...string) error

	// Query T from store.
	Query() Q
}

type BaseQuery[M any] interface {
	// Count all rows matching query.
	Count(ctx context.Context) (int, error)

	// Get all rows matching query.
	All(ctx context.Context) ([]*M, error)

	// Get first row matching query.
	One(ctx context.Context) (*M, error)

	// Delete all rows matching query.
	Delete(ctx context.Context) error
}
