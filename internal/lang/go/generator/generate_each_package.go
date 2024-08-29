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

func generateEachPackage(ctx context.Context, ormoptPackageImportPath string, packageName string, packageSource *source.PackageSource, tablesInPackage []*TableInfo) error {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	cfg := contexts.GenerateConfig(ctx)

	// each package
	eachPackageFileName := filepath.Base(strings.TrimSuffix(eachPackageTmpl, ".tmpl"))
	eachPackageFilePath := filepath.Join(cfg.GoORMOutputPath, packageSource.SourceRelativePath, eachPackageFileName)
	logs.Stdout.Debug("create file: file=" + eachPackageFilePath)
	var eachPackageFile *os.File
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.Create", func(_ context.Context) (err error) {
		eachPackageFile, err = os.Create(strings.TrimSuffix(eachPackageFilePath, ".go") + ".orm.gen.go")
		//nolint:wrapcheck
		return err
	}); err != nil {
		return errorz.Errorf("os.Create: %w", err)
	}
	defer eachPackageFile.Close()

	_ = tracez.StartFuncWithSpanNameSuffix(ctx, "template.New", func(_ context.Context) (err error) {
		eachPackageTemplateOnce.Do(func() {
			eachPackageTemplate = template.Must(template.New(eachPackageTmpl).Funcs(templateFuncMap(cfg)).Parse(string(mustz.One(templates.ReadFile(eachPackageTmpl)))))
		})
		return nil
	})

	eachPackageFileBuf := bytes.NewBuffer(nil)
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "eachPackageTemplate.Execute", func(_ context.Context) (err error) {
		return eachPackageTemplate.Execute(eachPackageFileBuf, FileInfo{
			SourceFile:              packageSource.SourceRelativePath,
			PackageName:             packageName,
			PackageImportPath:       packageSource.PackageImportPath,
			ORMOptPackageImportPath: ormoptPackageImportPath,
			Dialect:                 cfg.Dialect,
			SliceTypeSuffix:         cfg.GoSliceTypeSuffix,
			Tables:                  tablesInPackage,
		})
	}); err != nil {
		return errorz.Errorf("eachPackageTemplate.Execute: %w", err)
	}

	content := eachPackageFileBuf.Bytes()

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "eachPackageFile.Write", func(_ context.Context) (err error) {
		_, err = eachPackageFile.Write(content)
		//nolint:wrapcheck
		return err
	}); err != nil {
		return errorz.Errorf("eachPackageFile.Write: %w", err)
	}

	return nil
}
