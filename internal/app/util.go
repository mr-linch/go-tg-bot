package app

import (
	"github.com/rs/zerolog"
)

type tgbLoggerAdapater struct {
	logger zerolog.Logger
}

func (adapter tgbLoggerAdapater) Printf(format string, v ...interface{}) {
	adapter.logger.Printf(format, v...)
}
