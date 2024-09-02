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
}

func (s *StructSource) CommentGroupString() string {
	var builder strings.Builder
	for _, comment := range s.CommentGroup.List {
		builder.WriteString(comment.Text + "\n")
	}
	return builder.String()
}

func (s *StructSource) ExtractTableName(ctx context.Context, goColumnTag string) string {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	for _, comment := range s.CommentGroup.List {
		if matches := GoColumnTagCommentLineRegex(ctx, goColumnTag).FindStringSubmatch(comment.Text); len(matches) > _GoColumnTagCommentLineRegexTagValueIndex {
			return matches[_GoColumnTagCommentLineRegexTagValueIndex]
		}
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
		s.CommentGroupString(),
	)
}

const (
	//	                                              _____________ <- 1. comment prefix
	//	                                                              __ <- 2. tag key
	//	                                                                                      ___ <- 5. comment suffix
	_GoColumnTagCommentLineRegexFormat        = `^\s*(//+\s*|/\*\s*)?(%s)\s*:\s*(\S*)\s*(\S*)(\*/)?`
	_GoColumnTagCommentLineRegexTagNameIndex  = /*                               ^^^ 3. tag name */ 3
	_GoColumnTagCommentLineRegexTagValueIndex = /*                                       ^^^ 4. tag value */ 4
)

//nolint:gochecknoglobals
var (
	_GoColumnTagCommentLineRegexMap sync.Map
)

func GoColumnTagCommentLineRegex(ctx context.Context, goColumnTag string) *regexp.Regexp {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	if v, ok := _GoColumnTagCommentLineRegexMap.Load(goColumnTag); ok {
		return v.(*regexp.Regexp)
	}

	re := regexp.MustCompile(fmt.Sprintf(_GoColumnTagCommentLineRegexFormat, goColumnTag))
	_GoColumnTagCommentLineRegexMap.Store(goColumnTag, re)

	return re
}
