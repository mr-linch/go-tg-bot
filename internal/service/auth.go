package service

import (
	"context"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
)

//go:generate mockery --name Auth

type Auth interface {
	// AuthViaBot authorize user via bot.
	// If user is not exist, create new user.
	AuthViaBot(ctx context.Context, user *tg.User) (*domain.User, error)

	// AuthViaWidget authorize user via Telegram Login Widget.
	// AuthViaWidget(ctx context.Context, data *tg.AuthWidget) (*domain.User, error)

	// AuthViaWeb authorize user via WebAppInitData
	// AuthViaWebApp(ctx context.Context, data *tg.WebAppInitData) (*domain.User, error)
}
