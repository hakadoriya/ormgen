package source

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/pkg/apperr"
	"github.com/hakadoriya/z.go/buildz"
	"github.com/hakadoriya/z.go/contextz"
	"github.com/hakadoriya/z.go/errorz"
)

func walkDirFn(ctx context.Context, sourcePath string, fileExt string, packageSources *PackageSourceSlice) func(path string, d fs.DirEntry, err error) error {
	return func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return errorz.Errorf("path=%s: %w", filePath, err)
		}

		if err := contextz.CheckContext(ctx); err != nil {
			return errorz.Errorf("contextz.CheckContext: %w", err)
		}

		if d.IsDir() || !strings.HasSuffix(filePath, fileExt) || strings.HasSuffix(filePath, "_test.go") {
			contexts.Trace(ctx).Debug(fmt.Sprintf("skip: path=%s", filePath))
			return nil
		}

		cfg := contexts.GenerateConfig(ctx)
		packageImportPath := cfg.GoORMStructPackageImportPath
		if packageImportPath == "" {
			packageImportPath, err = buildz.FindPackageImportPath(sourcePath)
			if err != nil {
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
			PackageImportPath:  filepath.Join(packageImportPath, filepath.Dir(fileSource.SourceRelativePath)),
			SourceRelativePath: filepath.Dir(fileSource.SourceRelativePath),
			FileSources:        FileSourceSlice{fileSource},
		})

		return nil
	}
}
