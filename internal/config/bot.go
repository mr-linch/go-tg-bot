package config

type Bot struct {
	Token string `required:"true" usage:"telegram bot api token"`
}
