package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func GetLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}

func WithHandler(base *slog.Logger, name string) *slog.Logger {
	return base.With(slog.String("handler", name))
}

func WithService(base *slog.Logger, name string) *slog.Logger {
	return base.With(slog.String("service", name))
}

func WithRepo(base *slog.Logger, name string) *slog.Logger {
	return base.With(slog.String("repo", name))
}

func WithMethod(base *slog.Logger, name string) *slog.Logger {
	return base.With(slog.String("method", name))
}
