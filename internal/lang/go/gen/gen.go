package gen

import (
	"context"
	"embed"
	"go/ast"
	"os"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

type EachGoTmpl struct {
	SourceFile        string
	PackageName       string
	PackageImportPath string
	Structs           []Struct
}

type Struct struct {
	StructName   string
	TableName    string
	FieldsNames  []string
	FieldsTypes  []string
	ColumnsNames []string
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

			var structs []Struct
			for _, structSource := range fileSource.StructSources {
				tableName := structSource.ExtractTableName(ctx)
				var fields []string
				var FieldsTypes []string
				var columns []string
				for _, field := range structSource.StructType.Fields.List {
					fields = append(fields, field.Names[0].Name)
					FieldsTypes = append(FieldsTypes, fieldName(field.Type).String())
					columns = append(columns, reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get(cfg.GoColumnTag))
				}
				structs = append(structs, Struct{
					StructName:   structSource.TypeSpec.Name.Name,
					TableName:    tableName,
					FieldsNames:  fields,
					ColumnsNames: columns,
				})
			}

			template.Must(template.New("orm").Funcs(template.FuncMap{
				"add":        func(a, b int) int { return a + b },
				"UpperFirst": func(s string) string { return strings.ToUpper(string(s[0])) + s[1:] },
				"LowerFirst": func(s string) string { return strings.ToLower(string(s[0])) + s[1:] },
				"Base":       path.Base,
				"PlaceHolder": func(i int, column string) string {
					switch cfg.Dialect {
					case consts.DialectPostgres:
						return "$" + strconv.Itoa(i)
					}
					return "?"
				},
			}).Parse(string(mustz.One(eachGoTmpl.ReadFile(eachGoTmplName))))).Execute(f, EachGoTmpl{
				SourceFile:        fileSource.SourceRelativePath,
				PackageName:       packageSource.PackageName,
				PackageImportPath: packageSource.PackageImportPath,
				Structs:           structs,
			})

			defer f.Close()
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
