package source

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/otelz/tracez"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/iface"
	"github.com/hakadoriya/ormgen/pkg/apperr"
)

const fileExt = ".go"

func Parse(ctx context.Context, args []string) (PackageSourceSlice, error) {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	return parse(ctx, &iface.Pkg{FilepathAbsFunc: filepath.Abs}, args)
}

func parse(ctx context.Context, pkg iface.Iface, args []string) (PackageSourceSlice, error) {
	if len(args) != 1 {
		return nil, errorz.Errorf("invalid number of arguments; expected 1, got %d: %w", len(args), apperr.ErrInvalidArguments)
	}

	sourcePath, err := pkg.FilepathAbs(args[0])
	if err != nil {
		return nil, errorz.Errorf("pkg.FilepathAbs: %w", err)
	}

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return nil, errorz.Errorf("os.Stat: %w", err)
	}

	if !sourceInfo.IsDir() {
		return nil, errorz.Errorf("sourceInfo.IsDir: sourcePath=%s: %w", sourcePath, apperr.ErrSourcePathIsNotDirectory)
	}

	var packageSources PackageSourceSlice
	if err := filepath.WalkDir(sourcePath, newParseWalkDir(ctx, sourcePath, fileExt, &packageSources)); err != nil {
		return nil, errorz.Errorf("filepath.WalkDir: %w", err)
	}

	// DEBUG
	for _, packageSource := range packageSources {
		contexts.Trace(ctx).Debug("packageSource", slog.String("packageSource", fmt.Sprintf("%#v", packageSource)))
		for _, fileSource := range packageSource.FileSources {
			contexts.Trace(ctx).Debug("fileSource", slog.String("fileSource", fmt.Sprintf("%#v", fileSource)))
			for _, structSource := range fileSource.StructSources {
				contexts.Trace(ctx).Debug("structSource", slog.String("structSource", fmt.Sprintf("%#v", structSource)))
			}
		}
	}

	return packageSources, nil
}
