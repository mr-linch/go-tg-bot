package bot

import (
	"context"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

func (bot *Bot) registerMyChatMemberHandlers() {
	bot.MyChatMember(bot.onMyChatMemberPrivateCmd,
		tgb.ChatType(tg.ChatTypePrivate),
	)
}

func (bot *Bot) onMyChatMemberPrivateCmd(ctx context.Context, cm *tgb.ChatMemberUpdatedUpdate) error {
	me, err := cm.Client.Me(ctx)
	if err != nil {
		return err
	}

	if me.ID != cm.NewChatMember.User.ID ||
		me.ID != cm.OldChatMember.User.ID {
		return nil
	}

	user, err := bot.Service.User().AuthViaBot(ctx, &cm.From, nil)
	if err != nil {
		return err
	}

	if cm.OldChatMember.Status == "member" &&
		cm.NewChatMember.Status == "kicked" {
		return bot.Service.User().Stop(ctx, user)
	}

	if cm.OldChatMember.Status == "kicked" &&
		cm.NewChatMember.Status == "member" {
		return bot.Service.User().Restart(ctx, user)
	}

	return nil
}
