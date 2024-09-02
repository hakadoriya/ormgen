package source

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/util"
	"github.com/hakadoriya/ormgen/pkg/apperr"
	"github.com/hakadoriya/z.go/errorz"
)

const fileExt = ".go"

func Parse(ctx context.Context, args []string) (PackageSourceSlice, error) {
	if len(args) != 1 {
		return nil, errorz.Errorf("invalid number of arguments; expected 1, got %d: %w", len(args), apperr.ErrInvalidArguments)
	}

	sourcePath := util.Abs(ctx, args[0])

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return nil, errorz.Errorf("os.Stat: %w", err)
	}

	if !sourceInfo.IsDir() {
		return nil, errorz.Errorf("sourceInfo.IsDir: sourcePath=%s: %w", sourcePath, apperr.ErrSourcePathIsNotDirectory)
	}

	var packageSources PackageSourceSlice
	if err := filepath.WalkDir(sourcePath, walkDirFn(ctx, sourcePath, fileExt, &packageSources)); err != nil {
		return nil, errorz.Errorf("filepath.WalkDir: %w", err)
	}

	// DEBUG
	for _, packageSource := range packageSources {
		contexts.Stdout(ctx).Debug("packageSource", slog.String("packageSource", fmt.Sprintf("%#v", packageSource)))
		for _, fileSource := range packageSource.FileSources {
			contexts.Stdout(ctx).Debug("fileSource", slog.String("fileSource", fmt.Sprintf("%#v", fileSource)))
			for _, structSource := range fileSource.StructSources {
				contexts.Stdout(ctx).Debug("structSource", slog.String("structSource", fmt.Sprintf("%#v", structSource)))
			}
		}
	}

	return packageSources, nil
}
