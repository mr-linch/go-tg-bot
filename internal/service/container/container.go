package container

import (
	"github.com/mr-linch/go-tg-bot/internal/service"
	"github.com/mr-linch/go-tg-bot/internal/service/auth"
	"github.com/mr-linch/go-tg-bot/internal/store"
)

type Deps struct {
	Store store.Store
}

type Container struct {
	auth service.Auth
}

var _ service.Service = (*Container)(nil)

func New(deps Deps) *Container {
	return &Container{
		auth: &auth.Service{
			Store: deps.Store,
		},
	}
}

func (c *Container) Auth() service.Auth {
	return c.auth
}
