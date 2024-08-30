package logs

import (
	"io"
	"log/slog"
	"os"

	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/z.go/logz/slogz"
)

//nolint:gochecknoglobals
var (
	Trace  = New(io.Discard, slog.LevelDebug)
	Stdout = New(os.Stdout, slog.LevelInfo)
	Stderr = New(os.Stderr, slog.LevelDebug)
)

func New(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slogz.NewHandler(w, level)).With(slog.String("app", consts.AppName))
}
