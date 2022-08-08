package service

import (
	"context"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
)

//go:generate mockery --name Auth

type AuthSignUpOpts struct {
	Deeplink string
}

type Auth interface {
	// AuthViaBot authorize user via bot.
	// If user is not exist, create new user.
	AuthViaBot(ctx context.Context, user *tg.User, opts *AuthSignUpOpts) (*domain.User, error)

	// AuthViaWidget authorize user via Telegram Login Widget.
	// AuthViaWidget(ctx context.Context, data *tg.AuthWidget) (*domain.User, error)

	// AuthViaWeb authorize user via WebAppInitData
	// AuthViaWebApp(ctx context.Context, data *tg.WebAppInitData) (*domain.User, error)

	// Change language of user.
	SetUserLanguage(ctx context.Context, user *domain.User, lang string) (changed bool, err error)
}
