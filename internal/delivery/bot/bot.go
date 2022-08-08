package bot

import (
	"context"
	"regexp"

	"github.com/rs/zerolog"

	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/locales"
	"github.com/mr-linch/go-tg-bot/internal/service"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Deps struct {
	Bundle  *i18n.Bundle
	Logger  *zerolog.Logger
	Service service.Service
}

type Bot struct {
	*tgb.Router
	*Deps
}

func New(deps *Deps) (*Bot, error) {
	bot := &Bot{
		Deps:   deps,
		Router: tgb.NewRouter(),
	}

	bot.Use(tgb.MiddlewareFunc(func(next tgb.Handler) tgb.Handler {
		return tgb.HandlerFunc(func(ctx context.Context, update *tgb.Update) error {
			ctx = deps.Logger.WithContext(ctx)
			return next.Handle(ctx, update)
		})
	}))

	bot.registerHandlers()

	return bot, nil
}

func (bot *Bot) registerHandlers() *Bot {
	bot.registerLangHandlers()
	bot.registerGeneralHandlers()

	return bot
}

func (bot *Bot) T(user *domain.User) *i18n.Localizer {
	return i18n.NewLocalizer(
		bot.Bundle,
		user.PreferredLanguageCode.String,
		user.LanguageCode.String,
		locales.Default.String(),
	)
}

func getRegexpNamedArgs(re *regexp.Regexp, text string) map[string]string {
	match := re.FindStringSubmatch(text)

	result := make(map[string]string, len(match))

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result
}
