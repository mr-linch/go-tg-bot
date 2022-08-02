package container

import (
	"github.com/benbjohnson/clock"
	"github.com/mr-linch/go-tg-bot/internal/service"
	"github.com/mr-linch/go-tg-bot/internal/service/auth"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Deps struct {
	Store  store.Store
	Clock  clock.Clock
	Bundle *i18n.Bundle
}

type Container struct {
	auth service.Auth
}

var _ service.Service = (*Container)(nil)

func New(deps Deps) *Container {
	return &Container{
		auth: &auth.Service{
			Store:  deps.Store,
			Clock:  deps.Clock,
			Bundle: deps.Bundle,
		},
	}
}

func (c *Container) Auth() service.Auth {
	return c.auth
}
