package source

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/pkg/apperr"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/pathz/filepathz"
)

// MEMO: `sourcePath` only needs to calculate the relative path from the `sourcePath` to the `filePath`.
func parseFile(ctx context.Context, sourcePath, filePath string) (*FileSource, error) {
	contexts.Stdout(ctx).Debug(fmt.Sprintf("parse file: filename=%s", filePath))

	fset := token.NewFileSet()
	rootNode, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		// MEMO: parser.ParseFile err contains file path, so no need to log it
		return nil, errorz.Errorf("parser.ParseFile=%v: %w", err, apperr.ErrNoSourceFound)
	}

	sourceRelativePath, err := filepath.Rel(sourcePath, filePath)
	if err != nil {
		return nil, errorz.Errorf("filepath.Rel: %w", err)
	}

	structSources := make(StructSourceSlice, 0, 1)

	for commentedNode, commentGroups := range ast.NewCommentMap(fset, rootNode, rootNode.Comments) {
		for _, commentGroup := range commentGroups {
		CommentGroupLoop:
			for _, commentLine := range commentGroup.List {
				contexts.Trace(ctx).Debug(fmt.Sprintf("comment=%s: %s", filepathz.ExtractShortPath(fset.Position(commentGroup.Pos()).String()), commentLine.Text))
				// NOTE: If the comment line matches the GoColumnTag, it is assumed to be a comment line for the struct.
				cfg := contexts.GenerateConfig(ctx)
				if matches := GoColumnTagCommentLineRegex(ctx, cfg.GoColumnTag).FindStringSubmatch(commentLine.Text); len(matches) > _GoColumnTagCommentLineRegexTagNameIndex {
					ast.Inspect(commentedNode, func(node ast.Node) bool {
						switch nod := node.(type) {
						case *ast.TypeSpec:
							typeSpec := nod
							switch typ := typeSpec.Type.(type) {
							case *ast.StructType:
								structType := typ
								if structHasGoColumnTag(ctx, structType) {
									position := fset.Position(structType.Pos())
									contexts.Stdout(ctx).Debug(fmt.Sprintf("found struct source:%s: tag=%s:%s, type=%s", position, cfg.GoColumnTag, matches[_GoColumnTagCommentLineRegexTagNameIndex], typeSpec.Name.Name))
									structSources = append(structSources, &StructSource{
										PackageName:  rootNode.Name.Name,
										Position:     position,
										TypeSpec:     typeSpec,
										StructType:   structType,
										CommentGroup: commentGroup,
									})
								}
								return false
							default: // noop
								contexts.Stderr(ctx).Warn(fmt.Sprintf("unexpected type has tag: tag=%s:%s, type=%s", cfg.GoColumnTag, matches[_GoColumnTagCommentLineRegexTagNameIndex], typ))
							}
						default: // noop
						}
						return true
					})
					break CommentGroupLoop // NOTE: There may be multiple "GoColumnTag"s in the same commentGroup, so once you find the first one, break.
				}
			}
		}
	}

	if len(structSources) == 0 {
		return nil, errorz.Errorf("path=%s: %w", filePath, apperr.ErrNoSourceFound)
	}

	sort.SliceStable(structSources, func(i, j int) bool { return structSources[i].Position.Line < structSources[j].Position.Line })

	return &FileSource{
		FilePath:           filePath,
		PackageName:        rootNode.Name.Name,
		SourceRelativePath: sourceRelativePath,
		StructSources:      structSources,
	}, nil
}

func structHasGoColumnTag(ctx context.Context, s *ast.StructType) bool {
	cfg := contexts.GenerateConfig(ctx)

	for _, field := range s.Fields.List {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			if columnName := tag.Get(cfg.GoColumnTag); columnName != "" {
				return true
			}
		}
	}

	return false
}
