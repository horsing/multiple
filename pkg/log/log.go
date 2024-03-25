package log

import (
	"github.com/rs/zerolog/log"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

var loggerNameSeparator = "/"

func New(name string) *logr.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	zerologr.NameFieldName = name
	zerologr.NameSeparator = loggerNameSeparator
	zerologr.SetMaxV(1)

	l := zerologr.New(&log.Logger)
	return &l
}