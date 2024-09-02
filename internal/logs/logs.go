package logs

import (
	"io"
	"log/slog"
	"os"

	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/z.go/logz/slogz"
)

type TraceLogger struct{ *slog.Logger }
type StdoutLogger struct{ *slog.Logger }
type StderrLogger struct{ *slog.Logger }

//nolint:gochecknoglobals
var (
	Trace  *TraceLogger  = NewTrace(io.Discard, slog.LevelDebug)
	Stdout *StdoutLogger = NewStdout(os.Stdout, slog.LevelInfo)
	Stderr *StderrLogger = NewStderr(os.Stderr, slog.LevelDebug)
)

func New(w io.Writer, level slog.Level, attrs ...slog.Attr) *slog.Logger {
	return slog.New(slogz.NewHandler(w, level)).With(slog.String("app", consts.AppName))
}

func NewTrace(w io.Writer, level slog.Level, attrs ...slog.Attr) *TraceLogger {
	return &TraceLogger{New(w, level, append(attrs, slog.String("logName", "trace"))...)}
}

func NewStdout(w io.Writer, level slog.Level, attrs ...slog.Attr) *StdoutLogger {
	return &StdoutLogger{New(w, level, append(attrs, slog.String("logName", "stdout"))...)}
}

func NewStderr(w io.Writer, level slog.Level, attrs ...slog.Attr) *StderrLogger {
	return &StderrLogger{New(w, level, append(attrs, slog.String("logName", "stderr"))...)}
}
