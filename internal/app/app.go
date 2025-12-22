package app

import (
	"os"

	"github.com/rs/zerolog"
)

func SetupLogger(env string) zerolog.Logger {
	level := zerolog.InfoLevel
	if env == "development" {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)

	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
}
