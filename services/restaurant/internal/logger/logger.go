package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

type loggerContextKey struct{}

var key loggerContextKey

func NewLogger(config Config) *Logger {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: config.GetLevel(),
			},
		),
	)

	return &Logger{Logger: logger}
}

func ToContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(
		ctx,
		key,
		l,
	)
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}
