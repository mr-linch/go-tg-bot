package tgx

import "github.com/mr-linch/go-tg"

type TextMessage struct {
	text                  string
	replyMarkup           *tg.ReplyMarkup
	disableWebPagePreview *bool
	entities              *[]tg.MessageEntity
	parseMode             *tg.ParseMode
}

func NewTextMessage(text string) *TextMessage {
	return &TextMessage{
		text: text,
	}
}

func (msg *TextMessage) ReplyMarkup(rm tg.ReplyMarkup) *TextMessage {
	msg.replyMarkup = &rm
	return msg
}

func (msg *TextMessage) DisableWebPagePreview(b bool) *TextMessage {
	msg.disableWebPagePreview = &b
	return msg
}

func (msg *TextMessage) Entities(entities []tg.MessageEntity) *TextMessage {
	msg.entities = &entities
	return msg
}

func (msg *TextMessage) ParseMode(pm tg.ParseMode) *TextMessage {
	msg.parseMode = &pm
	return msg
}

func (msg *TextMessage) AsSendCall(peer tg.PeerID) *tg.SendMessageCall {
	call := tg.NewSendMessageCall(peer, msg.text)

	if msg.replyMarkup != nil {
		call.ReplyMarkup(*msg.replyMarkup)
	}

	if msg.disableWebPagePreview != nil {
		call.DisableWebPagePreview(*msg.disableWebPagePreview)
	}

	if msg.entities != nil {
		call.Entities(*msg.entities)
	}

	if msg.parseMode != nil {
		call.ParseMode(*msg.parseMode)
	}

	return call
}

func (msg *TextMessage) AsEditTextCall(peer tg.PeerID, msgID int) *tg.EditMessageTextCall {
	call := tg.NewEditMessageTextCall(peer, msgID, msg.text)

	if msg.replyMarkup != nil {
		if ikm, ok := (*msg.replyMarkup).(tg.InlineKeyboardMarkup); ok {
			call.ReplyMarkup(ikm)
		}
	}

	if msg.disableWebPagePreview != nil {
		call.DisableWebPagePreview(*msg.disableWebPagePreview)
	}

	if msg.entities != nil {
		call.Entities(*msg.entities)
	}

	if msg.parseMode != nil {
		call.ParseMode(*msg.parseMode)
	}

	return call
}

func (msg *TextMessage) AsEditTextInlineCall(inlineMsgID string) *tg.EditMessageTextCall {
	call := tg.NewEditMessageTextInlineCall(inlineMsgID, msg.text)

	if msg.replyMarkup != nil {
		if ikm, ok := (*msg.replyMarkup).(tg.InlineKeyboardMarkup); ok {
			call.ReplyMarkup(ikm)
		}
	}

	if msg.disableWebPagePreview != nil {
		call.DisableWebPagePreview(*msg.disableWebPagePreview)
	}

	if msg.entities != nil {
		call.Entities(*msg.entities)
	}

	if msg.parseMode != nil {
		call.ParseMode(*msg.parseMode)
	}

	return call
}

func (msg *TextMessage) AsEditReplyMarkupCall(peer tg.PeerID, msgID int) *tg.EditMessageReplyMarkupCall {
	call := tg.NewEditMessageReplyMarkupCall(peer, msgID)

	if msg.replyMarkup != nil {
		if ikm, ok := (*msg.replyMarkup).(tg.InlineKeyboardMarkup); ok {
			call.ReplyMarkup(ikm)
		}
	}

	return call
}

func (msg *TextMessage) AsEditReplyMarkupInlineCall(inlineMsgID string) *tg.EditMessageReplyMarkupCall {
	call := tg.NewEditMessageReplyMarkupInlineCall(inlineMsgID)

	if msg.replyMarkup != nil {
		if ikm, ok := (*msg.replyMarkup).(tg.InlineKeyboardMarkup); ok {
			call.ReplyMarkup(ikm)
		}
	}

	return call
}
