package domain

import (
	"time"

	"github.com/mr-linch/go-tg"
	"github.com/volatiletech/null/v8"
)

// UserID is a unique identifier for a user.
type UserID int64

type User struct {
	// Unique identifier for user in go-tg-bot.
	ID UserID

	// ID of Telegram user.
	TelegramID tg.UserID

	// First name
	FirstName string

	// Last name (optional)
	LastName null.String

	// TelegramUsername (optional)
	TelegramUsername null.String

	// Latest client language code (optional)
	LanguageCode null.String

	// Preferred language code.
	// If language code is in the list of supported languages, then it will be used.
	// Otherwise, the default language code will be used.
	PreferredLanguageCode null.String

	// Deeplink of the user.
	// Contains start parameters of the bot.
	Deeplink null.String

	// Time when user was created.
	CreatedAt time.Time

	// UpdatedAt is the time when user was last updated (optional).
	UpdatedAt null.Time
}

func (user *User) Name() string {
	if user.LastName.Valid {
		return user.FirstName + " " + user.LastName.String
	}
	return user.FirstName
}
