package contexts

import (
	"context"

	"github.com/hakadoriya/ormgen/internal/logs"
)

type ctxKeyStdout struct{}

func Stdout(ctx context.Context) *logs.StdoutLogger {
	if ctx != nil {
		if logger, ok := ctx.Value(ctxKeyStdout{}).((*logs.StdoutLogger)); ok {
			return logger
		}
	}

	return logs.Stdout
}

func WithStdout(ctx context.Context, logger *logs.StdoutLogger) context.Context {
	return context.WithValue(ctx, ctxKeyStdout{}, logger)
}

type ctxKeyStderr struct{}

func Stderr(ctx context.Context) *logs.StderrLogger {
	if ctx != nil {
		if logger, ok := ctx.Value(ctxKeyStderr{}).((*logs.StderrLogger)); ok {
			return logger
		}
	}

	return logs.Stderr
}

func WithStderr(ctx context.Context, logger *logs.StderrLogger) context.Context {
	return context.WithValue(ctx, ctxKeyStderr{}, logger)
}

type ctxKeyTrace struct{}

func Trace(ctx context.Context) *logs.TraceLogger {
	if ctx != nil {
		if logger, ok := ctx.Value(ctxKeyTrace{}).((*logs.TraceLogger)); ok {
			return logger
		}
	}

	return logs.Trace
}

func WithTrace(ctx context.Context, logger *logs.TraceLogger) context.Context {
	return context.WithValue(ctx, ctxKeyTrace{}, logger)
}
