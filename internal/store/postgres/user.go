package postgres

import (
	"context"
	"database/sql"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/mr-linch/go-tg-bot/internal/store/postgres/dal"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type userStore struct {
	baseStore
}

var _ store.User = (*userStore)(nil)

func (s *userStore) Add(ctx context.Context, user *domain.User) error {
	row := s.toRow(user)
	if err := s.insertOne(ctx, row); err != nil {
		return err
	}

	user.ID = domain.UserID(row.ID)

	return nil
}

func (s *userStore) Update(ctx context.Context, user *domain.User, fields ...string) error {
	row := s.toRow(user)

	if err := s.updateOne(ctx, row, store.ErrUserNotFound, fields...); err != nil {
		return err
	}

	return nil
}

func (s *userStore) Query() store.UserQuery {
	return &userQuery{
		store: s,
	}
}

func (s *userStore) toRow(user *domain.User) *dal.User {
	return &dal.User{
		ID:                    int(user.ID),
		TelegramID:            int64(user.TelegramID),
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		TelegramUsername:      user.TelegramUsername,
		LanguageCode:          user.LanguageCode,
		PreferredLanguageCode: user.PreferredLanguageCode,
		Deeplink:              user.Deeplink,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
	}
}

func (s *userStore) fromRow(row *dal.User) *domain.User {
	return &domain.User{
		ID:                    domain.UserID(row.ID),
		TelegramID:            tg.UserID(row.TelegramID),
		FirstName:             row.FirstName,
		LastName:              row.LastName,
		TelegramUsername:      row.TelegramUsername,
		LanguageCode:          row.LanguageCode,
		PreferredLanguageCode: row.PreferredLanguageCode,
		Deeplink:              row.Deeplink,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}

func (s *userStore) fromRows(rows dal.UserSlice) []*domain.User {
	users := make([]*domain.User, len(rows))
	for i, row := range rows {
		users[i] = s.fromRow(row)
	}
	return users
}

type userQuery struct {
	store *userStore
	mods  []qm.QueryMod
}

var _ store.UserQuery = (*userQuery)(nil)

func (q *userQuery) ID(ids ...domain.UserID) store.UserQuery {
	idsInt := make([]int, len(ids))
	for i, id := range ids {
		idsInt[i] = int(id)
	}
	q.mods = append(q.mods, dal.UserWhere.ID.IN(idsInt))
	return q
}

func (q *userQuery) TelegramID(ids ...tg.UserID) store.UserQuery {
	idsInt := make([]int64, len(ids))
	for i, id := range ids {
		idsInt[i] = int64(id)
	}
	q.mods = append(q.mods, dal.UserWhere.TelegramID.IN(idsInt))
	return q
}

func (q *userQuery) All(ctx context.Context) ([]*domain.User, error) {
	executor := q.store.getExecutor(ctx)

	users, err := dal.Users(q.mods...).All(ctx, executor)
	if err != nil {
		return nil, err
	}

	return q.store.fromRows(users), nil
}

func (q *userQuery) One(ctx context.Context) (*domain.User, error) {
	executor := q.store.getExecutor(ctx)

	user, err := dal.Users(q.mods...).One(ctx, executor)
	if err == sql.ErrNoRows {
		return nil, store.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return q.store.fromRow(user), nil
}

func (q *userQuery) Count(ctx context.Context) (int, error) {
	executor := q.store.getExecutor(ctx)

	v, err := dal.Users(q.mods...).Count(ctx, executor)
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

func (q *userQuery) Delete(ctx context.Context) error {
	executor := q.store.getExecutor(ctx)

	v, err := dal.Users(q.mods...).DeleteAll(ctx, executor)
	if err != nil {
		return err
	}

	if v == 0 {
		return store.ErrUserNotFound
	}

	return nil
}
