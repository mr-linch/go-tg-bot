package bot

import (
	"context"
	"fmt"
	"regexp"

	"github.com/friendsofgo/errors"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/locales"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
)

func (bot *Bot) buildLangKeyboard() tg.InlineKeyboardMarkup {
	languageTags := bot.Bundle.LanguageTags()

	buttons := make([]tg.InlineKeyboardButton, len(languageTags))
	for i, tag := range languageTags {
		info, ok := locales.Meta[tag.String()]
		if !ok {
			continue
		}

		buttons[i] = tg.NewInlineKeyboardButtonCallback(
			info.Emoji+" "+info.Label,
			fmt.Sprintf("lang_set:%s", tag.String()),
		)
	}

	return tg.NewInlineKeyboardMarkup(
		tg.NewButtonColumn(buttons...)...,
	)
}

func (bot *Bot) registerLangHandlers() {
	bot.Message(bot.onLangCmd,
		tgb.Command("lang"),
		tgb.ChatType(tg.ChatTypePrivate),
	)

	bot.CallbackQuery(bot.onLangSetCallback,
		tgb.Regexp(cbqLangSet),
	)
}

var cbqLangSet = regexp.MustCompile(`lang_set:(?P<lang>[a-zA-Z-_]{2,})`)

func (bot *Bot) onLangSetCallback(ctx context.Context, cbq *tgb.CallbackQueryUpdate) error {
	user, err := bot.Service.Auth().AuthViaBot(ctx, &cbq.From)
	if err != nil {
		return errors.Wrap(err, "auth via bot")
	}

	args := getRegexpNamedArgs(cbqLangSet, cbq.Data)

	lang, ok := args["lang"]
	if !ok {
		return errors.Errorf("lang not found in callback data: `%s`", cbq.Data)
	}

	changed, err := bot.Service.Auth().SetUserLanguage(ctx, user, lang)
	if err != nil {
		return errors.Wrap(err, "set user language")
	}

	var cbqAnswerText string

	if changed {
		cbqAnswerText = bot.T(user).MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "lang_set_cbq_changed",
				Other: "👌 Language has been changed",
			},
		})
	} else {
		cbqAnswerText = bot.T(user).MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "lang_set_cbq_is_current",
				Other: "🤷‍♂️ That language is your current language, please select another one",
			},
		})
	}

	if err := cbq.Update.Reply(ctx,
		cbq.AnswerText(cbqAnswerText, false),
	); err != nil {
		log.Ctx(ctx).Warn().
			Err(err).
			Str("query_id", cbq.ID).
			Msg("error while responding to callback query")
	}

	if changed {
		return cbq.Update.Reply(ctx, bot.buildStartMsg(user).AsEditTextCall(cbq.Message.Chat, cbq.Message.ID))
	}

	return nil
}

func (bot *Bot) onLangCmd(ctx context.Context, mu *tgb.MessageUpdate) error {
	user, err := bot.Service.Auth().AuthViaBot(ctx, mu.Message.From)
	if err != nil {
		return errors.Wrap(err, "auth via bot")
	}

	return mu.Answer(bot.T(user).MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "lang_message",
			Other: "Select your language from the list below:",
		},
	})).
		ParseMode(tg.HTML).
		ReplyMarkup(bot.buildLangKeyboard()).
		DoVoid(ctx)
}
