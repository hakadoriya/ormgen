package source

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"sync"

	"github.com/hakadoriya/z.go/otelz/tracez"
)

type PackageSourceSlice []*PackageSource

func (pss *PackageSourceSlice) AddPackageSource(packageSource *PackageSource) {
	// If there is already a package with the same name, add it there.
	for _, ps := range *pss {
		if ps.PackageName == packageSource.PackageName {
			ps.FileSources = append(ps.FileSources, packageSource.FileSources...)
			return
		}
	}

	// If it doesn't exist, add it as a new one.
	*pss = append(*pss, packageSource)
}

type PackageSource struct {
	PackageName        string
	PackageImportPath  string
	DirPath            string
	SourceRelativePath string

	FileSources FileSourceSlice
}

type FileSourceSlice []*FileSource

type FileSource struct {
	PackageName        string
	FilePath           string
	SourceRelativePath string

	StructSources StructSourceSlice
}

type StructSourceSlice []*StructSource

type StructSource struct {
	PackageName  string
	Position     token.Position
	TypeSpec     *ast.TypeSpec
	StructType   *ast.StructType
	CommentGroup *ast.CommentGroup

	extractedTableNameMap sync.Map
}

func CommentGroupString(commentGroup *ast.CommentGroup) string {
	var builder strings.Builder
	for _, comment := range commentGroup.List {
		builder.WriteString(comment.Text + "\n")
	}
	return builder.String()
}

func (s *StructSource) ExtractTableName(ctx context.Context, goColumnTag string) string {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	if v, ok := s.extractedTableNameMap.Load(s); ok {
		if result, ok := v.(string); ok {
			return result
		}
	}

	if matches := GoColumnTagCommentLineRegex(ctx, goColumnTag).FindStringSubmatch(CommentGroupString(s.CommentGroup)); len(matches) > _GoColumnTagCommentLineRegexTagValueIndex {
		result := matches[_GoColumnTagCommentLineRegexTagValueIndex]
		s.extractedTableNameMap.Store(s, result)
		return result
	}

	return ""
}

func (s *StructSource) GoString() string {
	return fmt.Sprintf(
		"&StructSource{PackageName: %#v, TokenPosition: %#v, TypeSpec: %#v, StructType: %#v, CommentGroup: %q}",
		s.PackageName,
		s.Position,
		s.TypeSpec,
		s.StructType,
		CommentGroupString(s.CommentGroup),
	)
}

const (
	// Regex for parsing the tag value from the comment line.
	//
	// e.g.:
	//
	//	// Table is table_name struct
	//	//
	//	//db:table table_name
	//	type User struct {
	//		UserID   int    `db:"user_id"`
	//		Username string `db:"username"`
	//	}
	//	                                                    __ <- 1. tag key == struct annotation key (e.g. db)
	_GoColumnTagCommentLineRegexFormat        = `(?m)^\s*//(%s):(\S*)\s+(\S+)(\*/)?`
	_GoColumnTagCommentLineRegexTagNameIndex  = /*               ^^^ 2. tag name (e.g. table) */ 2
	_GoColumnTagCommentLineRegexTagValueIndex = /*                       ^^^ 3. tag value (e.g. table_name) */ 3
)

//nolint:gochecknoglobals
var (
	_GoColumnTagCommentLineRegexMap sync.Map
)

func GoColumnTagCommentLineRegex(ctx context.Context, goColumnTag string) *regexp.Regexp {
	// NOTE: ctx needs to be used in the future.
	//nolint:ineffassign,staticcheck,wastedassign
	ctx, span := tracez.Start(ctx)
	defer span.End()

	if v, ok := _GoColumnTagCommentLineRegexMap.Load(goColumnTag); ok {
		re, ok := v.(*regexp.Regexp)
		if ok {
			return re
		}
	}

	re := regexp.MustCompile(fmt.Sprintf(_GoColumnTagCommentLineRegexFormat, goColumnTag))
	_GoColumnTagCommentLineRegexMap.Store(goColumnTag, re)

	return re
}
