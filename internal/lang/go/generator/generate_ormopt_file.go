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

//nolint:cyclop,funlen
func generateORMOptFile(ctx context.Context) (ormcommonPackageImportPath string, ormoptPackageImportPath string, err error) {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	cfg := contexts.GenerateConfig(ctx)
	if cfg.GoTableFileOnly {
		return "", "", nil
	}

	ormcommonDirPath := filepath.Join(cfg.GoORMOutputPath, filepath.Base(filepath.Dir(ormcommonFilename)))
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.MkdirAll", func(_ context.Context) (err error) {
		return os.MkdirAll(ormcommonDirPath, consts.Perm0o775)
	}); err != nil {
		return "", "", errorz.Errorf("os.MkdirAll: %w", err)
	}

	ormoptDirPath := filepath.Join(cfg.GoORMOutputPath, filepath.Base(filepath.Dir(ormoptFilename)))
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.MkdirAll", func(_ context.Context) (err error) {
		return os.MkdirAll(ormoptDirPath, consts.Perm0o775)
	}); err != nil {
		return "", "", errorz.Errorf("os.MkdirAll: %w", err)
	}

	ormgenPackageImportPath := cfg.GoORMOutputPackageImportPath
	if ormgenPackageImportPath == "" {
		if err := tracez.StartFuncWithSpanNameSuffix(ctx, "buildz.FindPackageImportPath", func(_ context.Context) (err error) {
			ormgenPackageImportPath, err = buildz.FindPackageImportPath(cfg.GoORMOutputPath)
			//nolint:wrapcheck
			return err
		}); err != nil {
			return "", "", errorz.Errorf("buildz.FindPackageImportPath: %w", err)
		}
	}

	// common file
	ormcommonFilePath := filepath.Join(ormcommonDirPath, filepath.Base(ormcommonFilename))
	logs.Stdout.Debug("create file: file=" + ormcommonFilePath)
	var ormcommonFile *os.File
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.Create", func(_ context.Context) (err error) {
		ormcommonFile, err = os.Create(strings.TrimSuffix(ormcommonFilePath, ".go") + ".orm.gen.go")
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", "", errorz.Errorf("os.Create: %w", err)
	}
	defer ormcommonFile.Close()

	var optcommonFileContent []byte
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "templates.ReadFile", func(_ context.Context) (err error) {
		optcommonFileContent, err = templates.ReadFile(ormcommonFilename)
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", "", errorz.Errorf("templates.ReadFile: %w", err)
	}

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "ormcommonFile.Write", func(_ context.Context) (err error) {
		_, err = ormcommonFile.Write(optcommonFileContent)
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", "", errorz.Errorf("ormcommonFile.Write: %w", err)
	}

	// opt file
	ormoptFilePath := filepath.Join(ormoptDirPath, filepath.Base(ormoptFilename))
	logs.Stdout.Debug("create file: file=" + ormoptFilePath)
	var ormoptFile *os.File
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.Create", func(_ context.Context) (err error) {
		ormoptFile, err = os.Create(strings.TrimSuffix(ormoptFilePath, ".go") + ".orm.gen.go")
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", "", errorz.Errorf("os.Create: %w", err)
	}
	defer ormoptFile.Close()

	var ormoptFileContent []byte
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "templates.ReadFile", func(_ context.Context) (err error) {
		ormoptFileContent, err = templates.ReadFile(ormoptFilename)
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", "", errorz.Errorf("templates.ReadFile: %w", err)
	}

	ormoptFileContent = append([]byte(consts.GeneratedFileHeader+"\n"), ormoptFileContent...)

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "ormoptFile.Write", func(_ context.Context) (err error) {
		_, err = ormoptFile.Write(ormoptFileContent)
		//nolint:wrapcheck
		return err
	}); err != nil {
		return "", "", errorz.Errorf("commonFile.Write: %w", err)
	}

	return filepath.Join(ormgenPackageImportPath, filepath.Base(ormcommonDirPath)), filepath.Join(ormgenPackageImportPath, filepath.Base(ormoptDirPath)), nil
}
