package logger

import (
	"os"

	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

func Init(cfg *config.Config) *Logger {
	zeroLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	switch cfg.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "err", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // log info and above by default
	}
	return &Logger{&zeroLogger}
}
