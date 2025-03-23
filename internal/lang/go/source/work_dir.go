package source

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hakadoriya/z.go/buildz"
	"github.com/hakadoriya/z.go/contextz"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/otelz/tracez"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/pkg/apperr"
)

//nolint:cyclop
func newParseWalkDir(ctx context.Context, sourcePath string, fileExt string, packageSources *PackageSourceSlice) func(path string, d fs.DirEntry, err error) error {
	walkDir := func(filePath string, d fs.DirEntry, err error) error {
		ctx, span := tracez.StartWithSpanNameSuffix(ctx, "walkDir")
		defer span.End()

		if err != nil {
			return errorz.Errorf("path=%s: %w", filePath, err)
		}

		if err := contextz.CheckContext(ctx); err != nil {
			return errorz.Errorf("contextz.CheckContext: %w", err)
		}

		if d.IsDir() || !strings.HasSuffix(filePath, fileExt) || strings.HasSuffix(filePath, "_test.go") {
			contexts.Trace(ctx).Debug("skip: path=" + filePath)
			return nil
		}

		cfg := contexts.GenerateConfig(ctx)
		ormStructPackageImportPath := cfg.GoORMStructPackageImportPath
		if ormStructPackageImportPath == "" && !cfg.GoTableFileOnly /* if table file only, don't need to find package import path */ {
			if err := tracez.StartFuncWithSpanNameSuffix(ctx, "buildz.FindPackageImportPath", func(_ context.Context) (err error) {
				ormStructPackageImportPath, err = buildz.FindPackageImportPath(sourcePath)
				//nolint:wrapcheck
				return err
			}); err != nil {
				return errorz.Errorf("buildz.FindPackageImportPath: %w", err)
			}
		}

		fileSource, err := parseFile(ctx, sourcePath, filePath)
		if err != nil {
			if errors.Is(err, apperr.ErrNoSourceFound) {
				contexts.Trace(ctx).Debug(fmt.Sprintf("skip: path=%s: %s", filePath, err))
				return nil
			}
			return errorz.Errorf("parseFile: %w", err)
		}

		packageSources.AddPackageSource(&PackageSource{
			PackageName:        fileSource.PackageName,
			DirPath:            filepath.Dir(fileSource.FilePath),
			PackageImportPath:  filepath.Join(ormStructPackageImportPath, filepath.Dir(fileSource.SourceRelativePath)),
			SourceRelativePath: filepath.Dir(fileSource.SourceRelativePath),
			FileSources:        FileSourceSlice{fileSource},
		})

		return nil
	}

	return walkDir
}
