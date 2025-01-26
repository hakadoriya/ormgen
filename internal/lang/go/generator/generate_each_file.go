package generator

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/mustz"
	"github.com/hakadoriya/z.go/otelz/tracez"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/ormgen/internal/logs"
)

func generateEachFile(ctx context.Context, ormoptPackageImportPath string, packageName string, packageSource *source.PackageSource, fileSource *source.FileSource, tablesInFile []*TableInfo) error {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	cfg := contexts.GenerateConfig(ctx)
	if cfg.GoTableFileOnly {
		return nil
	}

	// each file
	eachFilePath := filepath.Join(cfg.GoORMOutputPath, fileSource.SourceRelativePath)
	logs.Stdout.Debug("create file: file=" + eachFilePath)
	var eachFile *os.File
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.Create", func(_ context.Context) (err error) {
		eachFile, err = os.Create(strings.TrimSuffix(eachFilePath, ".go") + ".orm.gen.go")
		//nolint:wrapcheck
		return err
	}); err != nil {
		return errorz.Errorf("os.Create: %w", err)
	}
	defer eachFile.Close()

	_ = tracez.StartFuncWithSpanNameSuffix(ctx, "template.New", func(_ context.Context) (err error) {
		eachFileTemplateOnce.Do(func() {
			eachFileTemplate = template.Must(template.New(eachFileTmpl).Funcs(templateFuncMap(cfg)).Parse(string(mustz.One(templates.ReadFile(eachFileTmpl)))))
		})
		return nil
	})

	eachFileBuf := bytes.NewBuffer(nil)

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "eachFileTemplate.Execute", func(_ context.Context) (err error) {
		return eachFileTemplate.Execute(eachFileBuf, FileInfo{
			SourceFile:              fileSource.SourceRelativePath,
			PackageName:             packageName,
			PackageImportPath:       packageSource.PackageImportPath,
			ORMOptPackageImportPath: ormoptPackageImportPath,
			Dialect:                 cfg.Dialect,
			SliceTypeSuffix:         cfg.GoSliceTypeSuffix,
			Tables:                  tablesInFile,
		})
	}); err != nil {
		return errorz.Errorf("eachFileTemplate.Execute: %w", err)
	}

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "eachFile.Write", func(_ context.Context) (err error) {
		_, err = eachFile.Write(eachFileBuf.Bytes())
		return err //nolint:wrapcheck
	}); err != nil {
		return errorz.Errorf("eachFile.Write: %w", err)
	}

	return nil
}
