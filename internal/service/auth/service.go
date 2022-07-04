package auth

import (
	"context"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/volatiletech/null/v8"
)

type Service struct {
	Store store.Store
}

func (srv *Service) AuthViaBot(ctx context.Context, tgUser *tg.User) (*domain.User, error) {
	user, err := srv.Store.User().Query().TelegramID(tgUser.ID).One(ctx)
	if err == store.ErrUserNotFound {
		user = &domain.User{
			TelegramID:       tgUser.ID,
			FirstName:        tgUser.FirstName,
			LastName:         null.NewString(tgUser.LastName, tgUser.LastName != ""),
			LanguageCode:     null.NewString(tgUser.LanguageCode, tgUser.LanguageCode != ""),
			TelegramUsername: null.NewString(string(tgUser.Username), tgUser.Username != ""),
			CreatedAt:        time.Now(),
		}

		if err := srv.Store.User().Add(ctx, user); err != nil {
			return nil, err
		}

		return user, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "query user")
	}

	if err := srv.updateUserIfNeed(ctx, user, tgUser); err != nil {
		return nil, errors.Wrap(err, "update if need")
	}

	return user, nil
}

func (srv *Service) updateUserIfNeed(ctx context.Context, user *domain.User, tgUser *tg.User) error {
	var fields []string

	if user.FirstName != tgUser.FirstName {
		user.FirstName = tgUser.FirstName
		fields = append(fields, store.UserFields.FirstName)
	}

	if user.LastName.String != tgUser.LastName {
		user.LastName = null.NewString(tgUser.LastName, tgUser.LastName != "")
		fields = append(fields, store.UserFields.LastName)
	}

	if user.LanguageCode.String != tgUser.LanguageCode {
		user.LanguageCode = null.NewString(tgUser.LanguageCode, tgUser.LanguageCode != "")
		fields = append(fields, store.UserFields.LanguageCode)
	}

	if user.TelegramUsername.String != string(tgUser.Username) {
		user.TelegramUsername = null.NewString(string(tgUser.Username), tgUser.Username != "")
		fields = append(fields, store.UserFields.TelegramUsername)
	}

	if len(fields) > 0 {
		fields = append(fields, store.UserFields.UpdatedAt)
		user.UpdatedAt = null.TimeFrom(time.Now())
		if err := srv.Store.User().Update(ctx, user, fields...); err != nil {
			return err
		}
	}

	return nil
}
