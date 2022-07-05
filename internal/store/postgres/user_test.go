package postgres

import (
	"math/rand"
	"testing"
	"time"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/stretchr/testify/assert"
)

func newUser() *domain.User {
	return &domain.User{
		TelegramID: tg.UserID(rand.Int()),
		FirstName:  "FirstName",
		CreatedAt:  time.Now().Truncate(time.Second),
	}
}

func TestUserStore_Add(t *testing.T) {

	ctx, s := newPostgres(t)

	t.Run("OK", func(t *testing.T) {
		user := newUser()

		err := s.user.Add(ctx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	t.Run("Duplicate", func(t *testing.T) {
		user := newUser()

		err := s.user.Add(ctx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		err = s.user.Add(ctx, user)
		assert.Error(t, err, "add same user again should return error")
	})
}

func TestUserStore_Update(t *testing.T) {
	ctx, s := newPostgres(t)

	t.Run("OK", func(t *testing.T) {
		user := newUser()

		err := s.user.Add(ctx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		user.FirstName = "FirstName2"
		err = s.user.Update(ctx, user, store.UserFields.FirstName)
		assert.NoError(t, err)

		userFromDB, err := s.User().Query().ID(user.ID).One(ctx)
		assert.NoError(t, err)
		assert.Equal(t, user.FirstName, userFromDB.FirstName)
	})

	t.Run("NotFound", func(t *testing.T) {
		user := newUser()

		err := s.user.Update(ctx, user, store.UserFields.FirstName)
		assert.ErrorIs(t, err, store.ErrUserNotFound, "update non-existing user should return error")
	})

}

func TestUserQuery(t *testing.T) {
	ctx, s := newPostgres(t)

	users := make([]*domain.User, 5)

	for i := 0; i < len(users); i++ {
		user := newUser()
		err := s.user.Add(ctx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		users[i] = user
	}

	t.Run("One", func(t *testing.T) {
		user, err := s.User().Query().ID(users[0].ID).One(ctx)
		assert.NoError(t, err)
		assert.Equal(t, users[0].TelegramID, user.TelegramID)
	})

	t.Run("OneNotFound", func(t *testing.T) {
		user, err := s.User().Query().ID(0).One(ctx)
		assert.ErrorIs(t, err, store.ErrUserNotFound, "query non-existing user should return error")
		assert.Nil(t, user)
	})

	t.Run("All", func(t *testing.T) {
		usersFromDB, err := s.User().Query().All(ctx)
		assert.NoError(t, err)
		assert.Len(t, usersFromDB, len(users))
	})

	t.Run("Count", func(t *testing.T) {
		count, err := s.User().Query().Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, len(users), count)
	})

	t.Run("TelegramID", func(t *testing.T) {
		user, err := s.User().Query().TelegramID(users[0].TelegramID).One(ctx)
		assert.NoError(t, err)
		assert.Equal(t, users[0].ID, user.ID)
	})

	t.Run("Delete", func(t *testing.T) {
		err := s.User().Query().Delete(ctx)
		assert.NoError(t, err)

		count, err := s.User().Query().Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("DeleteNotFound", func(t *testing.T) {
		err := s.User().Query().Delete(ctx)
		assert.ErrorIs(t, err, store.ErrUserNotFound, "delete non-existing user should return error")
	})
}
