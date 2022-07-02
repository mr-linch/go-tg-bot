package config

type Sentry struct {
	DSN string `default:"" usage:"sentry dsn"`
}
