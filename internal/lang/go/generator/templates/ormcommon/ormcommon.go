package ormcommon

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
)

type QueryerContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// ErrNotUnique is returned when a query returns multiple rows.
var ErrNotUnique = errors.New("not unique")

// NOTE: noopLogger needs var declare for default use.
//
//nolint:gochecknoglobals
var noopLogger = slog.New(slog.NewJSONHandler(io.Discard, nil))

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return noopLogger
	}
	if logger, ok := ctx.Value((*slog.Logger)(nil)).(*slog.Logger); ok {
		return logger
	}
	return noopLogger
}

func LoggerWithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, (*slog.Logger)(nil), logger)
}
