package entrypoint

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/grafana/pyroscope-go"
	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/grpcz/grpclogz"
	"github.com/hakadoriya/z.go/logz/slogz"
	"github.com/hakadoriya/z.go/otelz"
	"github.com/hakadoriya/z.go/otelz/tracez"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/grpclog"

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

//nolint:funlen,cyclop
func Generate(c *cliz.Command, args []string) (err error) {
	ctx := c.Context()

	// trace
	ctx, span := tracez.Start(ctx)
	defer span.End()

	// profile
	if pyroscopeEndpoint := os.Getenv("PYROSCOPE_ENDPOINT"); pyroscopeEndpoint != "" {
		//nolint:exhaustruct
		p, err := pyroscope.Start(pyroscope.Config{
			ApplicationName: "ormgen",
			ServerAddress:   pyroscopeEndpoint,
		})
		if err != nil {
			err = errorz.Errorf("pyroscope.Start: %w", err)
			logs.Stderr.Debug(err.Error(), slogz.Error(err))
		} else {
			defer func() {
				if err := p.Stop(); err != nil {
					err = errorz.Errorf("p.Stop: %w", err)
					logs.Stderr.Debug(err.Error(), slogz.Error(err))
				}
			}()
		}
	}

	// log
	grpclog.SetLoggerV2(grpclogz.NewGRPCLoggerV2(logs.Stdout.Logger))

	// otel
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

	logs.Otel = otelslog.NewLogger(consts.AppName, otelslog.WithSource(true)).With(slog.String("logName", "otel"))
	commandNamesStr := strings.Join(c.GetExecutedCommandNames(), " ")
	logs.Otel.InfoContext(ctx, "🔆 start "+commandNamesStr, slog.Any("os.Args", os.Args))
	defer logs.Otel.InfoContext(ctx, "💤 end "+commandNamesStr, slog.Any("os.Args", os.Args))

	if err := otelruntime.Start(); err != nil {
		err = errorz.Errorf("runtime.Start: %w", err)
		logs.Stderr.Debug(err.Error(), slogz.Error(err))
	}

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
