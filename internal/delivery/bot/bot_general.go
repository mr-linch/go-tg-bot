package bot

import (
	"context"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/service"
	"github.com/mr-linch/go-tg-bot/pkg/tgx"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (bot *Bot) onStartCmd(ctx context.Context, msg *tgb.MessageUpdate) error {
	opts := &service.AuthSignUpOpts{}

	_, before, found := strings.Cut(msg.Text, " ")
	if found {
		opts.Deeplink = before
	}

	user, err := bot.Service.Auth().AuthViaBot(ctx, msg.Message.From, opts)
	if err != nil {
		return errors.Wrap(err, "auth via bot")
	}

	return msg.Update.Reply(ctx, bot.buildStartMsg(user).AsSendCall(msg.Chat.ID))
}

func (bot *Bot) buildStartMsg(user *domain.User) *tgx.TextMessage {
	text := bot.T(user).MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "start_message",
			Other: "Hi, {{.User.Name}}!\n" +
				"\n" +
				"<u>Here some info about you:</u>\n" +
				"├ <strong>ID:</strong> <code>{{.User.ID}}</code>\n" +
				"├ <strong>Telegram ID:</strong>  <code>{{.User.TelegramID}}</code>\n" +
				"├ <strong>Deeplink:</strong> <code>{{.User.Deeplink.String}}</code>\n" +
				"└ <strong>Language:</strong>  <code>{{.User.LanguageCode.String}}</code>\n",
		},
		TemplateData: map[string]any{
			"User": user,
		},
	})

	return tgx.NewTextMessage(text).ParseMode(tg.HTML)
}

func (bot *Bot) registerGeneralHandlers() {
	bot.Message(bot.onStartCmd,
		tgb.Command("start"),
		tgb.ChatType(tg.ChatTypePrivate),
	)
}
