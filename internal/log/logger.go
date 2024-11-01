package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct {
	zerolog.Logger
}

func (cl *Logger) Log(args ...interface{}) error {
	cl.Error().Stack().Err(args[1].(error)).Msg("")
	return nil
}

func New() *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := &Logger{Logger: zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()}
	return logger
}
