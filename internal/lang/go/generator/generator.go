package generator

import (
	"bytes"
	"context"
	"embed"
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/mustz"
)

const (
	eachFileTmpl    = "templates/each_file.go.tmpl"
	eachPackageTmpl = "templates/each_package.go.tmpl"
	queryOptionTmpl = "templates/query_option.go"
)

var (
	//go:embed templates
	templates embed.FS

	eachFileTemplate     *template.Template
	eachFileTemplateOnce sync.Once

	eachPackageTemplate     *template.Template
	eachPackageTemplateOnce sync.Once
)

type FileInfo struct {
	SourceFile        string
	PackageName       string
	PackageImportPath string
	SliceTypeSuffix   string
	Tables            []*TableInfo
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

func templateFuncMap(cfg *config.GenerateConfig) template.FuncMap {
	return template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"sub":        func(a, b int) int { return a - b },
		"upperFirst": func(s string) string { return strings.ToUpper(string(s[0])) + s[1:] },
		"lowerFirst": func(s string) string { return strings.ToLower(string(s[0])) + s[1:] },
		"basename":   path.Base,
		"PlaceHolder": func(columns []*ColumnInfo, startIndex int) string {
			var builder strings.Builder
			for i := range columns {
				if i != 0 {
					builder.WriteString(", ")
				}
				switch cfg.Dialect {
				case consts.DialectPostgres:
					builder.WriteString("$")
					builder.WriteString(strconv.Itoa(i + startIndex))
				default:
					builder.WriteString("?")
				}
			}
			return builder.String()
		},
		"PlaceHolderInWhere": func(columns []*ColumnInfo, op string, startIndex int) string {
			var builder strings.Builder
			for i := range columns {
				if i != 0 {
					builder.WriteString(" " + op + " ")
				}
				builder.WriteString(columns[i].ColumnName)
				builder.WriteString(" = ")
				switch cfg.Dialect {
				case consts.DialectPostgres:
					builder.WriteString("$")
					builder.WriteString(strconv.Itoa(i + startIndex))
				default:
					builder.WriteString("?")
				}
			}
			return builder.String()
		},
	}
}

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

		var tablesInPackage []*TableInfo

		for _, fileSource := range packageSource.FileSources {
			filePath := path.Join(cfg.GoORMOutputPath, fileSource.SourceRelativePath)
			f, err := os.Create(filePath)
			if err != nil {
				return errorz.Errorf("os.Create: %w", err)
			}

			var tables []*TableInfo
			for _, structSource := range fileSource.StructSources {
				tableInfo := BuildTableInfo(ctx, structSource)
				tables = append(tables, tableInfo)
				tablesInPackage = append(tablesInPackage, tableInfo)
			}

			sort.Slice(tables, func(i, j int) bool { return tables[i].SortKey < tables[j].SortKey })

			eachFileTemplateOnce.Do(func() {
				eachFileTemplate = template.Must(template.New("orm").Funcs(templateFuncMap(cfg)).Parse(string(mustz.One(templates.ReadFile(eachFileTmpl)))))
			})

			if err := eachFileTemplate.Execute(f, FileInfo{
				SourceFile:        fileSource.SourceRelativePath,
				PackageName:       packageSource.PackageName,
				PackageImportPath: packageSource.PackageImportPath,
				SliceTypeSuffix:   cfg.GoSliceTypeSuffix,
				Tables:            tables,
			}); err != nil {
				return errorz.Errorf("template.Execute: %w", err)
			}

			defer f.Close()
		}

		packageFilePath := path.Join(cfg.GoORMOutputPath, packageSource.SourceRelativePath, "ormgen.go")
		packageFile, err := os.Create(packageFilePath)
		if err != nil {
			return errorz.Errorf("os.Create: %w", err)
		}
		defer packageFile.Close()

		eachPackageTemplateOnce.Do(func() {
			eachPackageTemplate = template.Must(template.New("orm").Funcs(templateFuncMap(cfg)).Parse(string(mustz.One(templates.ReadFile(eachPackageTmpl)))))
		})

		if err := eachPackageTemplate.Execute(packageFile, FileInfo{
			SourceFile:        packageSource.SourceRelativePath,
			PackageName:       packageSource.PackageName,
			PackageImportPath: packageSource.PackageImportPath,
			SliceTypeSuffix:   cfg.GoSliceTypeSuffix,
			Tables:            tablesInPackage,
		}); err != nil {
			return errorz.Errorf("template.Execute: %w", err)
		}

		queryOptionFilePath := path.Join(cfg.GoORMOutputPath, packageSource.SourceRelativePath, filepath.Base(queryOptionTmpl))
		queryOptionFile, err := os.Create(queryOptionFilePath)
		if err != nil {
			return errorz.Errorf("os.Create: %w", err)
		}
		defer queryOptionFile.Close()

		r, err := templates.ReadFile(queryOptionTmpl)
		if err != nil {
			return errorz.Errorf("os.ReadFile: %w", err)
		}

		r = bytes.Replace(r, []byte("package templates"), []byte("// Code generated by ormgen; DO NOT EDIT.\npackage "+packageSource.PackageName), 1)

		if _, err := queryOptionFile.Write(r); err != nil {
			return errorz.Errorf("queryOptionFile.Write: %w", err)
		}

	}

	return nil
}

type RegexIndex struct {
	Regex *regexp.Regexp
	Index int
}

//nolint:gochecknoglobals
var (
	RegexIndexTableName = RegexIndex{
		Regex: regexp.MustCompile(`^\s*(//+\s*|/\*\s*)?\S+\s*:\s*table(s)?\s*[: ]\s*(\S+.*)$`),
		Index: 3, //nolint:mnd // 3 is the index of the table name
	}
)
