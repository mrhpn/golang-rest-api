package app

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(env string) zerolog.Logger {
	var writer io.Writer

	if env == "development" {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	} else {
		// production: write json to both Stdout (for docker/k8s logging) and a rotating log file

		// ensure log directory exists
		logDir := "logs"
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			_ = os.Mkdir(logDir, 0755)
		}

		rotatingWriter := &lumberjack.Logger{
			Filename:   filepath.Join(logDir, "app.log"),
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}
		writer = zerolog.MultiLevelWriter(os.Stdout, rotatingWriter)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	l := zerolog.New(writer).
		With().
		Timestamp().
		Logger()

	log.Logger = l
	zerolog.DefaultContextLogger = &l

	return l
}
