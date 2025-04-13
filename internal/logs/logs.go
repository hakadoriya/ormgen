package logs

import (
	"io"
	"log/slog"
	"os"

	"github.com/hakadoriya/z.go/logz/slogz"
	"github.com/hakadoriya/z.go/slicez"
)

type (
	TraceLogger  struct{ *slog.Logger }
	StdoutLogger struct{ *slog.Logger }
	StderrLogger struct{ *slog.Logger }
)

//nolint:gochecknoglobals
var (
	Otel   *slog.Logger
	Trace  = NewTrace(io.Discard, slog.LevelDebug)
	Stdout = NewStdout(os.Stdout, slog.LevelInfo)
	Stderr = NewStderr(os.Stderr, slog.LevelInfo)
)

func New(w io.Writer, level slog.Level, attrs ...slog.Attr) *slog.Logger {
	return slog.New(slogz.NewHandler(w, level)).With(append(slicez.Map(attrs, func(_ int, a slog.Attr) any { return a }), slog.String("app", "ormgen"))...)
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
