package source

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hakadoriya/ormgen/internal/apperr"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/ormgen/internal/util"
	"github.com/hakadoriya/z.go/contextz"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/pathz/filepathz"
)

func Parse(ctx context.Context, args []string) error {
	if len(args) != 1 {
		logs.Stderr.ErrorContext(ctx, fmt.Sprintf("invalid number of arguments; expected 1, got %d", len(args)), slog.Any("args", args))
		return errorz.Errorf("invalid number of arguments; expected 1, got %d: %w", len(args), apperr.ErrInvalidArguments)
	}

	sourcePath := util.Abs(args[0])

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return errorz.Errorf("os.Stat: %w", sourcePath, err)
	}

	if !sourceInfo.IsDir() {
		return errorz.Errorf("sourceInfo.IsDir: sourcePath=%s: %w", sourcePath, apperr.ErrSourcePathIsNotDirectory)
	}

	var packageSources PackageSourceSlice
	if err := filepath.WalkDir(sourcePath, walkDirFn(ctx, sourcePath, &packageSources)); err != nil {
		return errorz.Errorf("filepath.WalkDir: %w", err)
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

var fileExt = ".go"

func walkDirFn(ctx context.Context, sourcePath string, packageSources *PackageSourceSlice) func(path string, d fs.DirEntry, err error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errorz.Errorf("path=%s: %w", path, err)
		}

		if err := contextz.CheckContext(ctx); err != nil {
			return errorz.Errorf("contextz.CheckContext: %w", err)
		}

		if d.IsDir() || !strings.HasSuffix(path, fileExt) || strings.HasSuffix(path, "_test.go") {
			logs.Trace.Debug(fmt.Sprintf("skip: path=%s", path))
			return nil
		}

		fileSource, err := parseFile(ctx, sourcePath, path)
		if err != nil {
			if errors.Is(err, apperr.ErrNoStructSourceFound) {
				logs.Trace.Debug(fmt.Sprintf("skip: path=%s: %v", path, err))
				return nil
			}
			return errorz.Errorf("parseFile: %w", err)
		}

		cfg := contexts.GenerateConfig(ctx)
		packageImportPath := cfg.GoORMStructPackageImportPath
		if packageImportPath == "" {
			packageImportPath, err = util.DetectPackageImportPath(sourcePath)
			if err != nil {
				return errorz.Errorf("util.DetectPackageImportPath: %w", err)
			}
		}

		packageSources.AddPackageSource(&PackageSource{
			PackageName:        fileSource.PackageName,
			PackageImportPath:  packageImportPath,
			SourceRelativePath: filepath.Dir(fileSource.SourceRelativePath),
			FileSources:        FileSourceSlice{fileSource},
		})

		return nil
	}
}

// MEMO: `sourcePath` only needs to calculate the relative path from the `sourcePath` to the `filePath`.
func parseFile(ctx context.Context, sourcePath, filePath string) (*FileSource, error) {
	logs.Stdout.Debug(fmt.Sprintf("parse file: filename=%s", filePath))

	fset := token.NewFileSet()
	rootNode, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		// MEMO: parser.ParseFile err contains file path, so no need to log it
		return nil, errorz.Errorf("parser.ParseFile: %w", err)
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
				logs.Trace.Debug(fmt.Sprintf("comment=%s: %s", filepathz.ExtractShortPath(fset.Position(commentGroup.Pos()).String()), commentLine.Text))
				// NOTE: If the comment line matches the GoColumnTag, it is assumed to be a comment line for the struct.
				if matches := GoColumnTagCommentLineRegex(ctx).FindStringSubmatch(commentLine.Text); len(matches) > _GoColumnTagCommentLineRegexContentIndex {
					ast.Inspect(commentedNode, func(node ast.Node) bool {
						switch nod := node.(type) {
						case *ast.TypeSpec:
							typeSpec := nod
							switch typ := typeSpec.Type.(type) {
							case *ast.StructType:
								structType := typ
								if structHasGoColumnTag(ctx, structType) {
									position := fset.Position(structType.Pos())
									logs.Stdout.Debug(fmt.Sprintf("found struct source:%s: overwrite with comment group: type=%s", position.String(), typeSpec.Name.Name))
									structSources = append(structSources, &StructSource{
										PackageName:   rootNode.Name.Name,
										TokenPosition: position,
										TypeSpec:      typeSpec,
										StructType:    structType,
										CommentGroup:  commentGroup,
									})
								}
								return false
							default: // noop
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
		return nil, errorz.Errorf("path=%s: %w", filePath, apperr.ErrNoStructSourceFound)
	}

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
