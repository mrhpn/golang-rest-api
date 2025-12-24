package app

import (
	"io"
	"os"
	"path/filepath"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(logCfg *config.LogConfig, env string) zerolog.Logger {
	var writer io.Writer

	if env == "development" {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	} else {
		// production: write json to both Stdout (for docker/k8s logging) and a rotating log file

		// ensure log directory exists
		if _, err := os.Stat(logCfg.Path); os.IsNotExist(err) {
			_ = os.Mkdir(logCfg.Path, 0755) // MkdirAll ?
		}

		rotatingWriter := &lumberjack.Logger{
			Filename:   filepath.Join(logCfg.Path, "app.log"),
			MaxSize:    logCfg.MaxSizeMB,
			MaxBackups: logCfg.MaxBackupCount,
			MaxAge:     logCfg.MaxAgeDay,
			Compress:   logCfg.Compress,
		}
		writer = zerolog.MultiLevelWriter(os.Stdout, rotatingWriter)
	}

	level, err := zerolog.ParseLevel(logCfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	l := zerolog.New(writer).
		With().
		Timestamp().
		Logger()

	log.Logger = l
	zerolog.DefaultContextLogger = &l

	return l
}
