package apperr

import (
	"log/slog"
	"os"

	"github.com/hakadoriya/z.go/logz/slogz"
)

//nolint:gochecknoglobals
var (
	Log = slog.New(slogz.NewHandler(os.Stderr, slog.LevelDebug))
)
