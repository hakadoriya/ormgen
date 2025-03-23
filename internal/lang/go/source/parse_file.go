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
	"sync"

	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/otelz/tracez"
	"github.com/hakadoriya/z.go/pathz/filepathz"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/pkg/apperr"
)

// MEMO: `sourcePath` only needs to calculate the relative path from the `sourcePath` to the `filePath`.
//
//nolint:funlen,cyclop,gocognit
func parseFile(ctx context.Context, sourcePath, filePath string) (*FileSource, error) {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	contexts.Stdout(ctx).Debug("parse file: file=" + filePath)

	fset := token.NewFileSet()
	var rootNode *ast.File
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "parser.ParseFile", func(_ context.Context) (err error) {
		rootNode, err = parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		//nolint:wrapcheck
		return err
	}); err != nil {
		// MEMO: parser.ParseFile err contains file path, so no need to log it
		return nil, errorz.Errorf("parser.ParseFile: %w", fmt.Errorf("%w: %w", err, apperr.ErrNoSourceFound))
	}

	sourceRelativePath, err := filepath.Rel(sourcePath, filePath)
	if err != nil {
		return nil, errorz.Errorf("filepath.Rel: %w", err)
	}

	structSources := make(StructSourceSlice, 0, 1)

	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "ast.Inspect", func(_ context.Context) (err error) {
		for commentedNode, commentGroups := range ast.NewCommentMap(fset, rootNode, rootNode.Comments) {
			for _, commentGroup := range commentGroups {
				commentGroupString := CommentGroupString(commentGroup)
				contexts.Trace(ctx).Debug(fmt.Sprintf("comment=%s: %q", filepathz.ExtractShortPath(fset.Position(commentGroup.Pos()).String()), commentGroupString))
				// NOTE: If the comment line matches the GoColumnTag, it is assumed to be a comment line for the struct.
				cfg := contexts.GenerateConfig(ctx)
				if matches := GoColumnTagCommentLineRegex(ctx, cfg.GoColumnTag).FindStringSubmatch(commentGroupString); len(matches) > _GoColumnTagCommentLineRegexTagNameIndex {
					ast.Inspect(commentedNode, func(node ast.Node) bool {
						switch nod := node.(type) {
						case *ast.TypeSpec:
							typeSpec := nod
							position := fset.Position(typeSpec.Pos())
							switch typeSpecType := typeSpec.Type.(type) {
							case *ast.StructType:
								structType := typeSpecType
								if structHasGoColumnTag(ctx, structType) {
									contexts.Stdout(ctx).Debug(fmt.Sprintf("found struct source: file=%s, tag=%s:%s, type=%s", position, cfg.GoColumnTag, matches[_GoColumnTagCommentLineRegexTagNameIndex], typeSpec.Name.Name))
									structSources = append(structSources, &StructSource{
										PackageName:           rootNode.Name.Name,
										Position:              position,
										TypeSpec:              typeSpec,
										StructType:            structType,
										CommentGroup:          commentGroup,
										extractedTableNameMap: sync.Map{},
									})
								}
								return false
							default:
								err = errorz.Errorf("unexpected type has comment annotation: file=%s, tag=%s:%s, type=%s: %w", position, cfg.GoColumnTag, matches[_GoColumnTagCommentLineRegexTagNameIndex], typeSpecType, apperr.ErrInvalidAnnotation)
								return false
							}
						default:
							// noop
						}
						return true
					})
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}); err != nil {
		return nil, errorz.Errorf("ast.Inspect: %w", err)
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
