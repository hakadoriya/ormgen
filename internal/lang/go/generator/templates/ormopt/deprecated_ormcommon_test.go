package ormopt

import (
	"context"
	"log/slog"
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestLoggerFromContext(t *testing.T) {
	t.Parallel()

	t.Run("success,nil_context", func(t *testing.T) {
		t.Parallel()

		actual := LoggerFromContext(nil)
		requirez.Equal(t, noopLogger, actual)
	})

	t.Run("success,slog.Default", func(t *testing.T) {
		t.Parallel()

		actual := LoggerFromContext(LoggerWithContext(context.Background(), slog.Default()))
		requirez.Equal(t, slog.Default(), actual)
	})

	t.Run("success,no_logger", func(t *testing.T) {
		t.Parallel()

		actual := LoggerFromContext(LoggerWithContext(context.Background(), slog.Default()))
		requirez.Equal(t, slog.Default(), actual)
	})
}
