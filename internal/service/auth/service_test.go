package auth

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/mr-linch/go-tg-bot/internal/store/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/volatiletech/null/v8"
)

func TestService_AuthViaBot(t *testing.T) {
	t.Run("NewUser", func(t *testing.T) {
		clock := clock.NewMock()
		clock.Set(time.Now())

		userStoreQuery := mocks.NewUserQuery(t)
		userStoreQuery.On("TelegramID", tg.UserID(1234)).Return(userStoreQuery)
		userStoreQuery.On("One", mock.Anything).Return(nil, store.ErrUserNotFound)

		userStore := mocks.NewUser(t)
		userStore.On("Query").Return(userStoreQuery)
		userStore.On("Add", mock.Anything, &domain.User{
			TelegramID:            tg.UserID(1234),
			FirstName:             "John",
			CreatedAt:             clock.Now(),
			PreferredLanguageCode: null.StringFrom("uk"),
			LanguageCode:          null.StringFrom("uk"),
		}).Return(nil).Run(func(args mock.Arguments) {
			user := args.Get(1).(*domain.User)
			user.ID = 1
		})

		store := mocks.NewStore(t)
		store.On("User").Return(userStore)

		service := Service{Store: store, Clock: clock}

		user, err := service.AuthViaBot(context.Background(), &tg.User{
			ID:           1234,
			FirstName:    "John",
			LanguageCode: "uk",
		})

		assert.NoError(t, err)
		assert.Equal(t, domain.UserID(1), user.ID)
	})

	t.Run("Existing", func(t *testing.T) {
		clock := clock.NewMock()
		clock.Set(time.Now())

		userStoreQuery := mocks.NewUserQuery(t)
		userStoreQuery.On("TelegramID", tg.UserID(1234)).Return(userStoreQuery)
		userStoreQuery.On("One", mock.Anything).Return(&domain.User{
			ID:         1,
			TelegramID: tg.UserID(1234),
			FirstName:  "John",
		}, nil)

		userStore := mocks.NewUser(t)
		userStore.On("Query").Return(userStoreQuery)
		userStore.On("Update", mock.Anything, &domain.User{
			ID:               1,
			TelegramID:       tg.UserID(1234),
			FirstName:        "Alex",
			LastName:         null.StringFrom("Smith"),
			TelegramUsername: null.StringFrom("@smith"),
			LanguageCode:     null.StringFrom("en"),
			UpdatedAt:        null.TimeFrom(clock.Now()),
		}, store.UserFields.FirstName, store.UserFields.LastName, store.UserFields.LanguageCode, store.UserFields.TelegramUsername, store.UserFields.UpdatedAt).Return(nil)

		store := mocks.NewStore(t)
		store.On("User").Return(userStore)

		service := Service{Store: store, Clock: clock}

		user, err := service.AuthViaBot(context.Background(), &tg.User{
			ID:           1234,
			FirstName:    "Alex",
			LastName:     "Smith",
			Username:     "@smith",
			LanguageCode: "en",
		})

		assert.NoError(t, err)
		assert.Equal(t, domain.UserID(1), user.ID)
	})

}
