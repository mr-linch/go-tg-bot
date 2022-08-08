package user

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/friendsofgo/errors"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/locales"
	"github.com/mr-linch/go-tg-bot/internal/service"
	"github.com/mr-linch/go-tg-bot/internal/store"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
)

type Service struct {
	Clock  clock.Clock
	Store  store.Store
	Bundle *i18n.Bundle
}

func (srv *Service) AuthViaBot(ctx context.Context, tgUser *tg.User, opts *service.AuthSignUpOpts) (*domain.User, error) {
	if opts == nil {
		opts = &service.AuthSignUpOpts{}
	}

	user, err := srv.Store.User().Query().TelegramID(tgUser.ID).One(ctx)
	if err == store.ErrUserNotFound {
		user = &domain.User{
			TelegramID:            tgUser.ID,
			FirstName:             tgUser.FirstName,
			LastName:              null.NewString(tgUser.LastName, tgUser.LastName != ""),
			LanguageCode:          null.NewString(tgUser.LanguageCode, tgUser.LanguageCode != ""),
			PreferredLanguageCode: null.NewString(tgUser.LanguageCode, tgUser.LanguageCode != ""),
			TelegramUsername:      null.NewString(string(tgUser.Username), tgUser.Username != ""),
			Deeplink:              null.NewString(opts.Deeplink, opts.Deeplink != ""),
			CreatedAt:             srv.Clock.Now(),
		}

		if err := srv.Store.User().Add(ctx, user); err != nil {
			return nil, errors.Wrap(err, "add user to store")
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

func (srv *Service) SetUserLanguage(ctx context.Context, user *domain.User, lang string) (bool, error) {
	log.Ctx(ctx).Info().
		Int("user_id", int(user.ID)).
		Str("old_lang", user.PreferredLanguageCode.String).
		Str("new_lang", lang).
		Msg("Auth.SetUserLanguage")

	_, ok := locales.Meta[lang]
	if !ok {
		return false, errors.Errorf("invalid language: %s", lang)
	}

	if user.PreferredLanguageCode.String == lang {
		return false, nil
	}

	user.PreferredLanguageCode.SetValid(lang)
	user.UpdatedAt = null.TimeFrom(srv.Clock.Now())

	if err := srv.Store.User().Update(ctx, user,
		store.UserFields.PreferredLanguageCode,
		store.UserFields.UpdatedAt,
	); err != nil {
		return false, errors.Wrap(err, "update user")
	}

	return true, nil
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
		user.UpdatedAt = null.TimeFrom(srv.Clock.Now())
		if err := srv.Store.User().Update(ctx, user, fields...); err != nil {
			return err
		}
	}

	return nil
}

func (srv *Service) Stop(ctx context.Context, user *domain.User) error {
	log.Ctx(ctx).Info().
		Int("user_id", int(user.ID)).
		Msg("Auth.Stop")

	user.StoppedAt = null.TimeFrom(srv.Clock.Now())
	user.UpdatedAt = null.TimeFrom(srv.Clock.Now())

	if err := srv.Store.User().Update(ctx, user,
		store.UserFields.StoppedAt,
		store.UserFields.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "update user")
	}

	return nil
}

func (srv *Service) Restart(ctx context.Context, user *domain.User) error {
	log.Ctx(ctx).Info().
		Int("user_id", int(user.ID)).
		Msg("Auth.Restart")

	user.StoppedAt = null.Time{}
	user.UpdatedAt = null.TimeFrom(srv.Clock.Now())

	if err := srv.Store.User().Update(ctx, user,
		store.UserFields.StoppedAt,
		store.UserFields.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "update user")
	}

	return nil
}
