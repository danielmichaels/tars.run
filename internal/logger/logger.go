package logger

import (
	"github.com/danielmichaels/shortlink-go/internal/config"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

func NewLogger(name string, cfg *config.Conf) zerolog.Logger {
	logger := httplog.NewLogger(name, httplog.Options{
		JSON:     cfg.Server.LogJson,
		Concise:  cfg.Server.LogConcise,
		LogLevel: cfg.Server.LogLevel,
	})
	return logger
}
