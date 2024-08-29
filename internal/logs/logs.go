package logs

import (
	"io"
	"log/slog"
	"os"

	"github.com/hakadoriya/z.go/logz/slogz"
)

//nolint:gochecknoglobals
var (
	Trace  = slog.New(slogz.NewHandler(io.Discard, slog.LevelDebug))
	Stdout = slog.New(slogz.NewHandler(os.Stdout, slog.LevelInfo))
	Stderr = slog.New(slogz.NewHandler(os.Stderr, slog.LevelDebug))
)

func New(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slogz.NewHandler(w, level))
}
