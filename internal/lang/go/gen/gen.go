package gen

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path"
	"text/template"

	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/mustz"
)

const (
	fileGoTmplName = "file.go.tmpl"
)

var (
	//go:embed file.go.tmpl
	fileGoTmpl embed.FS
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

			template.Must(template.New("orm").Parse(string(mustz.One(fileGoTmpl.ReadFile(fileGoTmplName))))).Execute(f, struct {
				SourceFile  string
				PackageName string
			}{
				SourceFile:  fileSource.SourceRelativePath,
				PackageName: packageSource.PackageName,
			})

			defer f.Close()
		}
	}

	// DEBUG
	logs.Stdout.Info("packageSources", slog.Any("packageSources", packageSources))
	for _, packageSource := range packageSources {
		logs.Stdout.Info("packageSource", slog.String("packageSource", fmt.Sprintf("%#v", packageSource)))
		for _, fileSource := range packageSource.FileSources {
			logs.Stdout.Info("fileSource", slog.String("fileSource", fmt.Sprintf("%#v", fileSource)))
			for _, structSource := range fileSource.StructSources {
				logs.Stdout.Info("structSource", slog.String("structSource", fmt.Sprintf("%#v", structSource)))
			}
		}
	}

	return nil
}
