package entrypoint

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hakadoriya/ormgen/internal/apperr"
	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/ormgen/internal/util"
	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/contextz"
	"github.com/hakadoriya/z.go/errorz"
)

func Generate(c *cliz.Command, args []string) error {
	if len(args) != 1 {
		logs.Stderr.ErrorContext(c.Context(), fmt.Sprintf("invalid number of arguments; expected 1, got %d", len(args)), slog.Any("args", args))
		return errorz.Errorf("invalid number of arguments; expected 1, got %d: %w", len(args), apperr.ErrInvalidArguments)
	}

	sourcePath := util.Abs(args[0])

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return errorz.Errorf("os.Stat: %w", sourcePath, err)
	}

	if !sourceInfo.IsDir() {
		return errorz.Errorf("sourceInfo.IsDir: sourcePath=%s: %w", sourcePath, apperr.ErrSourcePathIsNotDirectory)
	}

	if err := filepath.WalkDir(sourcePath, walkDirFn(c.Context())); err != nil {
		return errorz.Errorf("filepath.WalkDir: %w", err)
	}

	return nil
}

var fileExt = ".go"

func walkDirFn(ctx context.Context) func(path string, d fs.DirEntry, err error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errorz.Errorf("path=%s: %w", path, err)
		}

		if err := contextz.CheckContext(ctx); err != nil {
			return errorz.Errorf("contextz.CheckContext: %w", err)
		}

		if d.IsDir() || !strings.HasSuffix(path, fileExt) || strings.HasSuffix(path, "_test.go") {
			logs.Trace.Debug(fmt.Sprintf("skip: path=%s", path))
			return nil
		}

		logs.Stdout.Info("walkDirFn", slog.Any("path", path))

		return nil
	}
}
