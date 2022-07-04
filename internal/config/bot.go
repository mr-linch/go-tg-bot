package config

type Bot struct {
	Token string `required:"true" usage:"telegram bot api token"`

	Webhook struct {
		Path    string `default:"/webhook" required:"true" usage:"telegram bot webhook path"`
		BaseURL string `usage:"telegram bot webhook base url (without path), if not provided will be running in long-polling mode"`
	}
}
