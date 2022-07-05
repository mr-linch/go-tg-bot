package store

import (
	"github.com/friendsofgo/errors"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/store/postgres/dal"
)

var ErrUserNotFound = errors.New("user not found")

var UserFields = dal.UserColumns

//go:generate mockery --name UserQuery  --case underscore

type UserQuery interface {
	BaseQuery[domain.User]

	ID(ids ...domain.UserID) UserQuery
	TelegramID(ids ...tg.UserID) UserQuery
}

//go:generate mockery --name User --case underscore

type User interface {
	Base[domain.User, UserQuery]
}
