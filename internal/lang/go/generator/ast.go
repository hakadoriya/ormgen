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

	if cfg.GoSliceTypeSuffix != "" {
		astFile.Decls = append(astFile.Decls,
			// import
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote("strings"),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote("fmt"),
						},
					},
				},
			},
		)
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

			//	func (s StructNameSlice) GetMapByPK() map[string]StructNameSlice {
			//		result := make(map[string]StructNameSlice)
			//		for _, v := range s {
			//			pkStringSlice := []string{fmt.Sprint(v.PK1), fmt.Sprint(v.PK2), ...}
			//			pk := strings.Join(pkStringSlice, PKConcatenationSeparator)
			//			if _, ok := result[pk]; ok {
			//				panic(fmt.Errorf("duplicate primary key: table=%T, pk=%s, one=%v, other=%v", v, pk, v, result[pk]))
			//			}
			//			result[pk] = v
			//		}
			//		return result
			//	}
			pkStringSliceElts := make([]ast.Expr, 0)
			for _, c := range table.PrimaryKeys {
				// fmt.Sprint(v.PK1)
				pkStringSliceElts = append(pkStringSliceElts, &ast.CallExpr{
					Fun: &ast.SelectorExpr{X: &ast.Ident{Name: "fmt"}, Sel: &ast.Ident{Name: "Sprint"}},
					Args: []ast.Expr{
						// v.PK1
						&ast.BasicLit{Kind: token.STRING, Value: "v." + c.FieldName},
					},
				})
			}
			astFile.Decls = append(astFile.Decls,
				&ast.FuncDecl{
					Recv: &ast.FieldList{List: []*ast.Field{{
						Names: []*ast.Ident{{Name: "s"}},
						Type:  &ast.Ident{Name: structSource.TypeSpec.Name.Name + cfg.GoSliceTypeSuffix},
					}}},
					Name: &ast.Ident{Name: cfg.GoMapByPKMethod},
					Type: &ast.FuncType{
						Params:  &ast.FieldList{},
						Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "map[string]*" + structSource.TypeSpec.Name.Name}}}},
					},
					Body: &ast.BlockStmt{List: []ast.Stmt{
						// result := make(map[string]StructNameSlice)
						&ast.AssignStmt{
							Lhs: []ast.Expr{&ast.Ident{Name: "result"}},
							Tok: token.DEFINE, // :=
							Rhs: []ast.Expr{&ast.CallExpr{
								Fun:  &ast.Ident{Name: "make"},
								Args: []ast.Expr{&ast.Ident{Name: "map[string]*" + structSource.TypeSpec.Name.Name}},
							}},
						},
						// for _, v := range s {
						&ast.RangeStmt{
							Key:   &ast.Ident{Name: "_"},
							Value: &ast.Ident{Name: "v"},
							Tok:   token.DEFINE, // :=
							X:     &ast.Ident{Name: "s"},
							Body: &ast.BlockStmt{List: []ast.Stmt{
								//	pkStringSlice := []string{
								//		fmt.Sprintf("%v", v.PK1),
								//		fmt.Sprintf("%v", v.PK2),
								//		...
								//	}
								&ast.AssignStmt{
									Lhs: []ast.Expr{&ast.Ident{Name: "pkStringSlice"}},
									Tok: token.DEFINE, // :=
									Rhs: []ast.Expr{&ast.CompositeLit{Type: &ast.ArrayType{Elt: &ast.Ident{Name: "string"}}, Elts: pkStringSliceElts}},
								},
								// pk := strings.Join(pkStringSlice, PKConcatenationSeparator)
								&ast.AssignStmt{
									Lhs: []ast.Expr{&ast.Ident{Name: "pk"}},
									Tok: token.DEFINE, // :=
									Rhs: []ast.Expr{&ast.CallExpr{
										Fun:  &ast.Ident{Name: "strings.Join"},
										Args: []ast.Expr{&ast.Ident{Name: "pkStringSlice"}, &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(cfg.GoPKConcatenationSeparator)}},
									}},
								},
								// if _, ok := result[pk]; ok {
								&ast.IfStmt{
									Init: &ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{Name: "_"},
											&ast.Ident{Name: "ok"},
										},
										Tok: token.DEFINE, // :=
										Rhs: []ast.Expr{&ast.Ident{Name: "result[pk]"}},
									},
									Cond: &ast.Ident{Name: "ok"},
									Body: &ast.BlockStmt{List: []ast.Stmt{
										// panic(fmt.Errorf("duplicate primary key: pk=%s, one=%v, other=%v", pk, v, result[pk]))
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.Ident{Name: "panic"},
												Args: []ast.Expr{&ast.CallExpr{
													Fun:  &ast.Ident{Name: "fmt.Errorf"},
													Args: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote("duplicate primary key: table=%T, pk=%s, one=%v, other=%v")}, &ast.Ident{Name: "v"}, &ast.Ident{Name: "pk"}, &ast.Ident{Name: "v"}, &ast.Ident{Name: "result[pk]"}},
												}},
											},
										},
									}},
								},
								// result[pk] = v
								&ast.AssignStmt{
									Lhs: []ast.Expr{&ast.Ident{Name: "result[pk]"}},
									Tok: token.ASSIGN, // =
									Rhs: []ast.Expr{&ast.Ident{Name: "v"}},
								},
							}},
						},
						// return result
						&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "result"}}},
					}},
				},
			)
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
