package generator

import (
	"bytes"
	"context"
	"go/ast"
	"go/printer"
	"go/token"
	"io"
	"strconv"

	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/otelz/tracez"

	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
)

//nolint:exhaustruct,funlen
func fprintTableMethods(ctx context.Context, w io.Writer, fileSource *source.FileSource) error {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	cfg := contexts.GenerateConfig(ctx)

	astFile := &ast.File{
		// package
		Name: &ast.Ident{
			Name: fileSource.PackageName,
		},
		// methods
		Decls: []ast.Decl{},
	}

	for _, structSource := range fileSource.StructSources {
		//	func (s *StructName) TableName() string {
		//		return "TableName"
		//	}
		astFile.Decls = append(astFile.Decls,
			&ast.FuncDecl{
				Recv: &ast.FieldList{List: []*ast.Field{{
					Names: []*ast.Ident{{Name: "s"}},
					Type:  &ast.StarExpr{X: &ast.Ident{Name: structSource.TypeSpec.Name.Name}},
				}}},
				Name: &ast.Ident{Name: cfg.GoTableNameMethod},
				Type: &ast.FuncType{
					Params:  &ast.FieldList{},
					Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "string"}}}},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: strconv.Quote(structSource.ExtractTableName(ctx, cfg.GoColumnTag))}}}},
				},
			},
		)

		//
		// COLUMN
		//
		table := BuildTableInfo(ctx, structSource)
		// all column names method
		elts := make([]ast.Expr, 0)
		for _, c := range table.Columns {
			elts = append(elts, &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(c.ColumnName),
			})
		}
		astFile.Decls = append(astFile.Decls,
			//	func (s *StructName) Columns() []string {
			//		return []string{"column1", "column2", ...}
			//	}
			&ast.FuncDecl{
				Recv: &ast.FieldList{List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "s"}},
						Type:  &ast.StarExpr{X: &ast.Ident{Name: structSource.TypeSpec.Name.Name}},
					},
				}},
				Name: &ast.Ident{Name: cfg.GoColumnsNameMethod},
				Type: &ast.FuncType{
					Params:  &ast.FieldList{},
					Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "[]string"}}}},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.CompositeLit{Type: &ast.ArrayType{Elt: &ast.Ident{Name: "string"}}, Elts: elts}}}}},
			},
		)

		// each column name methods
		for i := range table.Columns {
			astFile.Decls = append(astFile.Decls,
				//	func (s *StructName) Column1() string {
				//		return "column1"
				//	}
				&ast.FuncDecl{
					Recv: &ast.FieldList{List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "s"}},
							Type:  &ast.StarExpr{X: &ast.Ident{Name: structSource.TypeSpec.Name.Name}},
						},
					}},
					Name: &ast.Ident{Name: cfg.GoColumnNameMethodPrefix + table.Columns[i].FieldName},
					Type: &ast.FuncType{
						Params:  &ast.FieldList{},
						Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "string"}}}},
					},
					Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: strconv.Quote(table.Columns[i].ColumnName)}}}}},
				},
			)
		}

		if cfg.GoSliceTypeSuffix != "" {
			astFile.Decls = append(astFile.Decls,
				// type StructNameSlice []*StructName
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{&ast.TypeSpec{
						Name: &ast.Ident{Name: structSource.TypeSpec.Name.Name + cfg.GoSliceTypeSuffix},
						Type: &ast.ArrayType{Elt: &ast.StarExpr{X: &ast.Ident{Name: structSource.TypeSpec.Name.Name}}},
					}},
				},
				//	func (s StructNameSlice) TableName() string {
				//		return "TableName"
				//	}
				&ast.FuncDecl{
					Recv: &ast.FieldList{List: []*ast.Field{{
						Names: []*ast.Ident{{Name: "s"}},
						Type:  &ast.Ident{Name: structSource.TypeSpec.Name.Name + cfg.GoSliceTypeSuffix},
					}}},
					Name: &ast.Ident{Name: cfg.GoTableNameMethod},
					Type: &ast.FuncType{
						Params:  &ast.FieldList{},
						Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "string"}}}},
					},
					Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: strconv.Quote(structSource.ExtractTableName(ctx, cfg.GoColumnTag))}}}}},
				},
			)

			//	func (s StructNameSlice) Columns() []string {
			//		return []string{"column1", "column2", ...}
			//	}
			astFile.Decls = append(astFile.Decls,
				&ast.FuncDecl{
					Recv: &ast.FieldList{List: []*ast.Field{{
						Names: []*ast.Ident{{Name: "s"}},
						Type:  &ast.Ident{Name: structSource.TypeSpec.Name.Name + cfg.GoSliceTypeSuffix},
					}}},
					Name: &ast.Ident{Name: cfg.GoColumnsNameMethod},
					Type: &ast.FuncType{
						Params:  &ast.FieldList{},
						Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "[]string"}}}},
					},
					Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.CompositeLit{Type: &ast.ArrayType{Elt: &ast.Ident{Name: "string"}}, Elts: elts}}}}},
				},
			)

			// each column name methods
			for i := range table.Columns {
				astFile.Decls = append(astFile.Decls,
					//	func (s StructNameSlice) Column1() string {
					//		return "column1"
					//	}
					&ast.FuncDecl{
						Recv: &ast.FieldList{List: []*ast.Field{{
							Names: []*ast.Ident{{Name: "s"}},
							Type:  &ast.Ident{Name: structSource.TypeSpec.Name.Name + cfg.GoSliceTypeSuffix},
						}}},
						Name: &ast.Ident{Name: cfg.GoColumnNameMethodPrefix + table.Columns[i].FieldName},
						Type: &ast.FuncType{
							Params:  &ast.FieldList{},
							Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "string"}}}},
						},
						Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: strconv.Quote(table.Columns[i].ColumnName)}}}}},
					},
				)
			}
		}
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.WriteString(
		consts.GeneratedFileHeader + "\n" +
			"//\n" +
			"// source: " + fileSource.SourceRelativePath + "\n" +
			"\n",
	); err != nil {
		return errorz.Errorf("buf.WriteString: %w", err)
	}

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "printer.Fprint", func(_ context.Context) (err error) {
		return printer.Fprint(buf, token.NewFileSet(), astFile)
	}); err != nil {
		return errorz.Errorf("printer.Fprint: %w", err)
	}

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "w.Write", func(_ context.Context) (err error) {
		content := bytes.ReplaceAll(buf.Bytes(), []byte("\n}\nfunc "), []byte("\n}\n\nfunc "))
		_, err = w.Write(content)
		//nolint:wrapcheck
		return err
	}); err != nil {
		return errorz.Errorf("w.Write: %w", err)
	}

	return nil
}
