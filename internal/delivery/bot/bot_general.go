package bot

import (
	"context"

	"github.com/friendsofgo/errors"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (bot *Bot) onStartCmd(ctx context.Context, msg *tgb.MessageUpdate) error {
	user, err := bot.Service.Auth().AuthViaBot(ctx, msg.Message.From)
	if err != nil {
		return errors.Wrap(err, "auth via bot")
	}

	return msg.Update.Respond(ctx, bot.buildStartMsg(user, msg.Chat))
}

func (bot *Bot) buildStartMsg(user *domain.User, peer tg.PeerID) *tg.SendMessageCall {
	return tg.NewSendMessageCall(peer, bot.T(user).MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "start_message",
			Other: "Hi! I'm @{{.Bot.Username}}\n" +
				"\n" +
				"Here some info about you:\n" +
				"├ <strong>ID:</strong> <code>{{.User.ID}}</code>\n" +
				"├ <strong>Telegram ID:</strong>  <code>{{.User.TelegramID}}</code>\n" +
				"└ <strong>Language:</strong>  <code>{{.User.LanguageCode.String}}</code>\n",
		},
		TemplateData: map[string]any{
			"Bot": struct {
				Username string
			}{
				Username: "remove_me",
			},
			"User": user,
		},
	})).ParseMode(tg.HTML)
}

func (bot *Bot) registerGeneralHandlers() {
	bot.Message(bot.onStartCmd,
		tgb.Command("start"),
		tgb.ChatType(tg.ChatTypePrivate),
	)
}
