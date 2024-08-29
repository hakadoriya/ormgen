package generator

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/hakadoriya/z.go/buildz"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/otelz/tracez"

	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/logs"
)

func generateORMOptFile(ctx context.Context) (ormoptPackageImportPath string, err error) {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	cfg := contexts.GenerateConfig(ctx)

	ormoptDirPath := filepath.Join(cfg.GoORMOutputPath, filepath.Base(filepath.Dir(ormoptTmpl)))
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.MkdirAll", func(_ context.Context) (err error) {
		return os.MkdirAll(ormoptDirPath, consts.Perm0o775)
	}); err != nil {
		return "", errorz.Errorf("os.MkdirAll: %w", err)
	}

	ormgenPackageImportPath := cfg.GoORMOutputPackageImportPath
	if ormgenPackageImportPath == "" {
		if err := tracez.StartFuncWithSpanNameSuffix(ctx, "buildz.FindPackageImportPath", func(_ context.Context) (err error) {
			ormgenPackageImportPath, err = buildz.FindPackageImportPath(cfg.GoORMOutputPath)
			//nolint:wrapcheck
			return err
		}); err != nil {
			return "", errorz.Errorf("buildz.FindPackageImportPath: %w", err)
		}
	}

	// common file
	ormoptFilePath := filepath.Join(ormoptDirPath, filepath.Base(ormoptTmpl))
	logs.Stdout.Debug("create file: file=" + ormoptFilePath)
	var ormoptFile *os.File
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.Create", func(_ context.Context) (err error) {
		ormoptFile, err = os.Create(strings.TrimSuffix(ormoptFilePath, ".go") + ".orm.gen.go")
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", errorz.Errorf("os.Create: %w", err)
	}
	defer ormoptFile.Close()

	var fileContent []byte
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "templates.ReadFile", func(_ context.Context) (err error) {
		fileContent, err = templates.ReadFile(ormoptTmpl)
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", errorz.Errorf("templates.ReadFile: %w", err)
	}

	fileContent = append([]byte(consts.GeneratedFileHeader+"\n"), fileContent...)

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "ormoptFile.Write", func(_ context.Context) (err error) {
		_, err = ormoptFile.Write(fileContent)
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", errorz.Errorf("commonFile.Write: %w", err)
	}

	return filepath.Join(ormgenPackageImportPath, filepath.Base(ormoptDirPath)), nil
}
