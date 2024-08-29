package entrypoint

import (
	"context"
	"log/slog"
	"runtime"

	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/grpcz/grpclogz"
	"github.com/hakadoriya/z.go/logz/slogz"
	"github.com/hakadoriya/z.go/otelz"
	"github.com/hakadoriya/z.go/otelz/tracez"
	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	gogenerator "github.com/hakadoriya/ormgen/internal/lang/go/generator"
	gosource "github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/ormgen/pkg/apperr"
)

//nolint:gochecknoglobals
var (
	memStats runtime.MemStats
)

func ki(bytes uint64) uint64 {
	const ki = 1024
	return bytes / ki
}

func Generate(c *cliz.Command, args []string) error {
	ctx := c.Context()

	grpclogz.NewGRPCLoggerV2(logs.Stdout.Logger)

	shutdown, err := otelz.SetupAutoExport(ctx, otelz.WithResourceOptions(resource.WithAttributes(attribute.String("service.name", consts.AppName))))
	if err != nil {
		err = errorz.Errorf("otelz.SetupAutoExport: %w", err)
		logs.Stderr.Debug(err.Error(), slogz.Error(err))
	}
	otel.SetErrorHandler(otelz.ErrorHandleFunc(func(err error) { logs.Stderr.Debug(err.Error(), slogz.Error(err)) }))
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), consts.GracefulShutdownTimeout)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			err = errorz.Errorf("shutdown: %w", err)
			logs.Stderr.Debug(err.Error(), slogz.Error(err))
		}
	}()

	if err := otelruntime.Start(); err != nil {
		err = errorz.Errorf("runtime.Start: %w", err)
		logs.Stderr.Debug(err.Error(), slogz.Error(err))
	}

	ctx, span := tracez.Start(ctx)
	defer span.End()

	runtime.ReadMemStats(&memStats)
	logs.Stdout.Debug("memStats", slog.Uint64("allocKi", ki(memStats.Alloc)), slog.Uint64("totalAllocKi", ki(memStats.TotalAlloc)), slog.Uint64("sysKi", ki(memStats.Sys)))
	defer logs.Stdout.Debug("memStats", slog.Uint64("allocKi", ki(memStats.Alloc)), slog.Uint64("totalAllocKi", ki(memStats.TotalAlloc)), slog.Uint64("sysKi", ki(memStats.Sys)))

	cfg := contexts.GenerateConfig(ctx)

	switch cfg.Language {
	case consts.LanguageGo:
		packageSources, err := gosource.Parse(ctx, args)
		if err != nil {
			return errorz.Errorf("gosource.Parse: %w", err)
		}

		if err := gogenerator.Generate(ctx, packageSources); err != nil {
			return errorz.Errorf("gogenerator.Generate: %w", err)
		}

	default:
		return errorz.Errorf("lang=%s: %w", cfg.Language, apperr.ErrLanguageNotSupported)
	}

	return nil
}
