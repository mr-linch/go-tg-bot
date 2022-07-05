package container

import (
	"testing"

	"github.com/mr-linch/go-tg-bot/internal/store/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	container := New(Deps{
		Store: mocks.NewStore(t),
	})

	assert.NotNil(t, container)
	assert.NotNil(t, container.Auth())
}
