package bot

import (
	"context"
	"fmt"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/service"
	"github.com/mr-linch/go-tg/tgb"
)

type Deps struct {
	Service service.Service
}

type Bot struct {
	*tgb.Router

	*Deps
}

func New(deps *Deps) *Bot {
	bot := &Bot{
		Deps:   deps,
		Router: tgb.NewRouter(),
	}

	bot.register()

	return bot
}

func (bot *Bot) register() *Bot {
	bot.Message(bot.onStart,
		tgb.Command("start"),
		tgb.ChatType(tg.ChatTypePrivate),
	)

	return bot
}

func (bot *Bot) onStart(ctx context.Context, msg *tgb.MessageUpdate) error {
	user, err := bot.Service.Auth().AuthViaBot(ctx, msg.Message.From)
	if err != nil {
		return err
	}

	return msg.Answer(tg.HTML.Text(
		fmt.Sprintf("Hey, %s!", user.FirstName),
		"",
		fmt.Sprintf("<strong>Your Bot ID:</strong> <code>%d</code>", user.ID),
		fmt.Sprintf("<strong>Your Telegram ID:</strong> <code>%d</code>", user.TelegramID),
	)).ParseMode(tg.HTML).DoVoid(ctx)
}
