package generator

import (
	"context"
	"embed"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/otelz/tracez"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/ormgen/pkg/apperr"
)

const (
	eachFileTmpl    = "templates/package/each_file.go.tmpl"
	eachPackageTmpl = "templates/package/ormgen.go.tmpl"
	ormoptTmpl      = "templates/ormopt/ormopt.go"
)

//nolint:gochecknoglobals
var (
	//go:embed templates
	templates embed.FS

	eachFileTemplate     *template.Template
	eachFileTemplateOnce sync.Once

	eachPackageTemplate     *template.Template
	eachPackageTemplateOnce sync.Once
)

type FileInfo struct {
	SourceFile              string
	PackageName             string
	PackageImportPath       string
	ORMOptPackageImportPath string
	Dialect                 string
	SliceTypeSuffix         string
	Tables                  []*TableInfo
}

func fieldName(x ast.Expr) *ast.Ident {
	switch t := x.(type) {
	case *ast.Ident:
		return t
	case *ast.SelectorExpr:
		if _, ok := t.X.(*ast.Ident); ok {
			return t.Sel
		}
	case *ast.StarExpr:
		return fieldName(t.X)
	}
	return nil
}

//nolint:cyclop
func templateFuncMap(cfg *config.GenerateConfig) template.FuncMap {
	return template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"sub":        func(a, b int) int { return a - b },
		"upperFirst": func(s string) string { return strings.ToUpper(string(s[0])) + s[1:] },
		"lowerFirst": func(s string) string { return strings.ToLower(string(s[0])) + s[1:] },
		"basename":   filepath.Base,
		"placeholder": func(columns []*ColumnInfo, startIndex int) string {
			var builder strings.Builder
			for i := range columns {
				if i != 0 {
					builder.WriteString(", ")
				}
				switch cfg.Dialect {
				case consts.DialectPostgres, consts.DialectCockroach:
					builder.WriteString("$")
					builder.WriteString(strconv.Itoa(i + startIndex))
				case consts.DialectMySQL, consts.DialectSQLite3, consts.DialectSpanner:
					builder.WriteString("?")
				default:
					panic(fmt.Errorf("dialect=%s: %w", cfg.Dialect, apperr.ErrDialectNotSupported))
				}
			}
			return builder.String()
		},
		"placeholderInWhere": func(columns []*ColumnInfo, op string, startIndex int) string {
			var builder strings.Builder
			for i := range columns {
				if i != 0 {
					builder.WriteString(" " + op + " ")
				}
				builder.WriteString(columns[i].ColumnName)
				builder.WriteString(" = ")
				switch cfg.Dialect {
				case consts.DialectPostgres, consts.DialectCockroach:
					builder.WriteString("$")
					builder.WriteString(strconv.Itoa(i + startIndex))
				case consts.DialectMySQL, consts.DialectSQLite3, consts.DialectSpanner:
					builder.WriteString("?")
				default:
					panic(fmt.Errorf("dialect=%s: %w", cfg.Dialect, apperr.ErrDialectNotSupported))
				}
			}
			return builder.String()
		},
	}
}

func Generate(ctx context.Context, packageSources source.PackageSourceSlice) error {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	cfg := contexts.GenerateConfig(ctx)

	ormoptPackageImportPath, err := generateORMOptFile(ctx)
	if err != nil {
		return errorz.Errorf("generateORMOptFile: %w", err)
	}

	for _, packageSource := range packageSources {
		if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.MkdirAll", func(_ context.Context) (err error) {
			packageDirPath := filepath.Join(cfg.GoORMOutputPath, packageSource.SourceRelativePath)
			return os.MkdirAll(packageDirPath, consts.Perm0o775)
		}); err != nil {
			return errorz.Errorf("os.MkdirAll: %w", err)
		}

		packageName := packageSource.PackageName
		if cfg.GoORMPackageName != "" {
			packageName = cfg.GoORMPackageName
		}

		var tablesInPackage []*TableInfo
		for _, fileSource := range packageSource.FileSources {
			if err := generateTableFile(ctx, fileSource); err != nil {
				return errorz.Errorf("generateTableFile: %w", err)
			}

			var tablesInFile []*TableInfo
			for _, structSource := range fileSource.StructSources {
				tableInfo := BuildTableInfo(ctx, structSource)
				tablesInFile = append(tablesInFile, tableInfo)
				tablesInPackage = append(tablesInPackage, tableInfo)
			}

			if err := generateEachFile(ctx, ormoptPackageImportPath, packageName, packageSource, fileSource, tablesInFile); err != nil {
				return errorz.Errorf("generateEachFile: %w", err)
			}
		}

		if err := generateEachPackage(ctx, ormoptPackageImportPath, packageName, packageSource, tablesInPackage); err != nil {
			return errorz.Errorf("generateEachPackage: %w", err)
		}
	}

	return nil
}
