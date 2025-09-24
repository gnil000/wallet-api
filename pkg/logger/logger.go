package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(level zerolog.Level) zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Level(level)
}

func WithModule(base zerolog.Logger, module string) zerolog.Logger {
	return base.With().Str("module", module).Logger()
}
