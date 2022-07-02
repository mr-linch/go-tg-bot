package config

type Postgres struct {
	DSN          string `required:"true" usage:"postgres dsn"`
	MaxOpenConns int    `default:"10" usage:"postgres max open conns"`
	MaxIdleConns int    `default:"0" usage:"postgres max idle conns"`
}
