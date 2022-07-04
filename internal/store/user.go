package store

import (
	"github.com/friendsofgo/errors"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/store/postgres/dal"
)

var ErrUserNotFound = errors.New("user not found")

var UserFields = dal.UserColumns

type UserQuery interface {
	BaseQuery[domain.User]

	ID(ids ...domain.UserID) UserQuery
	TelegramID(ids ...tg.UserID) UserQuery
}

type User interface {
	Base[domain.User, UserQuery]
}
