package gen

import (
	"context"
	"embed"
	"os"
	"path"
	"text/template"

	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/mustz"
)

const (
	eachGoTmplName = "each.go.tmpl"
)

var (
	//go:embed each.go.tmpl
	eachGoTmpl embed.FS
)

func Output(ctx context.Context, packageSources source.PackageSourceSlice) error {
	cfg := contexts.GenerateConfig(ctx)

	if err := os.MkdirAll(cfg.GoORMOutputPath, consts.Perm0o775); err != nil {
		return errorz.Errorf("os.MkdirAll: %w", err)
	}

	for _, packageSource := range packageSources {
		packageDirPath := path.Join(cfg.GoORMOutputPath, packageSource.SourceRelativePath)
		if err := os.MkdirAll(packageDirPath, consts.Perm0o775); err != nil {
			return errorz.Errorf("os.MkdirAll: %w", err)
		}

		for _, fileSource := range packageSource.FileSources {
			filePath := path.Join(cfg.GoORMOutputPath, fileSource.SourceRelativePath)
			f, err := os.Create(filePath)
			if err != nil {
				return errorz.Errorf("os.Create: %w", err)
			}

			template.Must(template.New("orm").Parse(string(mustz.One(eachGoTmpl.ReadFile(eachGoTmplName))))).Execute(f, struct {
				SourceFile        string
				PackageName       string
				PackageImportPath string
			}{
				SourceFile:        fileSource.SourceRelativePath,
				PackageName:       packageSource.PackageName,
				PackageImportPath: packageSource.PackageImportPath,
			})

			defer f.Close()
		}
	}

	return nil
}
