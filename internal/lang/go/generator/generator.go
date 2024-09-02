package generator

import (
	"context"
	"embed"
	"go/ast"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/ormgen/internal/util"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/mustz"
)

const (
	eachFileTmpl    = "templates/package/each_file.go.tmpl"
	eachPackageTmpl = "templates/package/ormgen.go.tmpl"
	commonTmpl      = "templates/ormgen/common.go"
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
	SourceFile              string
	PackageName             string
	PackageImportPath       string
	CommonPackageImportPath string
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

func templateFuncMap(cfg *config.GenerateConfig) template.FuncMap {
	return template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"sub":        func(a, b int) int { return a - b },
		"upperFirst": func(s string) string { return strings.ToUpper(string(s[0])) + s[1:] },
		"lowerFirst": func(s string) string { return strings.ToLower(string(s[0])) + s[1:] },
		"basename":   filepath.Base,
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

	commonDirPath := filepath.Join(cfg.GoORMOutputPath, "ormgen")
	if err := os.MkdirAll(commonDirPath, consts.Perm0o775); err != nil {
		return errorz.Errorf("os.MkdirAll: %w", err)
	}

	commonPackageImportPath, err := util.DetectPackageImportPath(commonDirPath)
	if err != nil {
		return errorz.Errorf("util.DetectPackageImportPath: %w", err)
	}

	commonFilePath := filepath.Join(commonDirPath, filepath.Base(commonTmpl))
	commonFile, err := os.Create(commonFilePath)
	if err != nil {
		return errorz.Errorf("os.Create: %w", err)
	}
	defer commonFile.Close()

	r, err := templates.ReadFile(commonTmpl)
	if err != nil {
		return errorz.Errorf("templates.ReadFile: %w", err)
	}

	if _, err := commonFile.Write(append([]byte(consts.GeneratedFileHeader+"\n"), r...)); err != nil {
		return errorz.Errorf("commonFile.Write: %w", err)
	}

	for _, packageSource := range packageSources {
		packageDirPath := filepath.Join(cfg.GoORMOutputPath, packageSource.SourceRelativePath)
		if err := os.MkdirAll(packageDirPath, consts.Perm0o775); err != nil {
			return errorz.Errorf("os.MkdirAll: %w", err)
		}

		var tablesInPackage []*TableInfo

		for _, fileSource := range packageSource.FileSources {
			eachFilePath := filepath.Join(cfg.GoORMOutputPath, fileSource.SourceRelativePath)
			eachFile, err := os.Create(eachFilePath)
			if err != nil {
				return errorz.Errorf("os.Create: %w", err)
			}

			var tables []*TableInfo
			for _, structSource := range fileSource.StructSources {
				tableInfo := BuildTableInfo(ctx, structSource)
				tables = append(tables, tableInfo)
				tablesInPackage = append(tablesInPackage, tableInfo)
			}

			eachFileTemplateOnce.Do(func() {
				eachFileTemplate = template.Must(template.New(eachFileTmpl).Funcs(templateFuncMap(cfg)).Parse(string(mustz.One(templates.ReadFile(eachFileTmpl)))))
			})

			if err := eachFileTemplate.Execute(eachFile, FileInfo{
				SourceFile:              fileSource.SourceRelativePath,
				PackageName:             packageSource.PackageName,
				PackageImportPath:       packageSource.PackageImportPath,
				CommonPackageImportPath: commonPackageImportPath,
				Dialect:                 cfg.Dialect,
				SliceTypeSuffix:         cfg.GoSliceTypeSuffix,
				Tables:                  tables,
			}); err != nil {
				return errorz.Errorf("template.Execute: %w", err)
			}

			defer eachFile.Close()
		}

		eachPackageFileName := filepath.Base(strings.TrimSuffix(eachPackageTmpl, ".tmpl"))
		eachPackageFilePath := filepath.Join(cfg.GoORMOutputPath, packageSource.SourceRelativePath, eachPackageFileName)
		eachPackageFile, err := os.Create(eachPackageFilePath)
		if err != nil {
			return errorz.Errorf("os.Create: %w", err)
		}
		defer eachPackageFile.Close()

		eachPackageTemplateOnce.Do(func() {
			eachPackageTemplate = template.Must(template.New(eachPackageTmpl).Funcs(templateFuncMap(cfg)).Parse(string(mustz.One(templates.ReadFile(eachPackageTmpl)))))
		})

		if err := eachPackageTemplate.Execute(eachPackageFile, FileInfo{
			SourceFile:              packageSource.SourceRelativePath,
			PackageName:             packageSource.PackageName,
			PackageImportPath:       packageSource.PackageImportPath,
			CommonPackageImportPath: commonPackageImportPath,
			SliceTypeSuffix:         cfg.GoSliceTypeSuffix,
			Tables:                  tablesInPackage,
		}); err != nil {
			return errorz.Errorf("template.Execute: %w", err)
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
