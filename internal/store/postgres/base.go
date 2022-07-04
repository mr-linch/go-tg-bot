package postgres

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/mr-linch/go-tg-bot/internal/store/postgres/shared"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// baseStore define base store
type baseStore struct {
	*sql.DB
	store.Txier
}

func (bs *baseStore) getExecutor(ctx context.Context) boil.ContextExecutor {
	return shared.GetExecutorOrDefault(ctx, bs)
}

// type deletableRow interface {
// 	Delete(
// 		ctx context.Context,
// 		exec boil.ContextExecutor,
// 	) (int64, error)
// }

// func (bs *BaseStore) deleteOne(
// 	ctx context.Context,
// 	row deletableRow,
// 	notFoundErr error,
// ) error {
// 	return bs.Txier(ctx, func(ctx context.Context) error {
// 		count, err := row.Delete(
// 			ctx,
// 			bs.getExecutor(ctx),
// 		)

// 		if err != nil {
// 			return errors.Wrap(err, "exec")
// 		}

// 		switch {
// 		case count == 0:
// 			return notFoundErr
// 		case count > 1:
// 			return store.ErrTooManyAffectedRows
// 		}

// 		return nil
// 	})
// }

type updatetableRow interface {
	Update(
		ctx context.Context,
		exec boil.ContextExecutor,
		columns boil.Columns,
	) (int64, error)
}

type insertableRow interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

func (bs *baseStore) insertOne(
	ctx context.Context,
	row insertableRow,
) error {
	return row.Insert(ctx, bs.getExecutor(ctx), boil.Infer())
}

func (bs *baseStore) updateOne(
	ctx context.Context,
	row updatetableRow,
	notFoundErr error,
	fields ...string,
) error {
	return bs.Txier(ctx, func(ctx context.Context) error {
		columns := boil.Infer()

		if len(fields) > 0 {
			columns = boil.Whitelist(fields...)
		}

		count, err := row.Update(
			ctx,
			bs.getExecutor(ctx),
			columns,
		)

		if err != nil {
			return errors.Wrap(err, "exec")
		}

		switch {
		case count == 0:
			return notFoundErr
		case count > 1:
			return store.ErrTooManyAffectedRows
		}

		return nil
	})
}

// func isPostgresErrorConstraint(err error, constraint string) bool {
// 	var pqErr *pq.Error
// 	if errors.As(err, &pqErr) {
// 		return pqErr.Constraint == constraint
// 	}
// 	return false
// }
