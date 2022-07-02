package config

type Log struct {
	Pretty bool `default:"false" usage:"pretty print log"`
	Debug  bool `default:"false" usage:"debug log level"`
}
