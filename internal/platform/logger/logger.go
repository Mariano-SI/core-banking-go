package logger

import (
	"io"
	"log/slog"

	"github.com/Mariano-SI/core-banking-go/internal/platform/config"
)

func New(w io.Writer, cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}

	var handler slog.Handler
	if cfg.IsProduction() {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	return slog.New(handler)
}
