package service

import (
	"context"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
)

type Auth interface {
	Auth(ctx context.Context, user *tg.User) (*domain.User, error)
}
