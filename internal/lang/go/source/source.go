package source

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"sync"

	"github.com/hakadoriya/ormgen/internal/contexts"
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
	SourceRelativePath string

	FileSources FileSourceSlice
}

type FileSourceSlice []*FileSource

type FileSource struct {
	FilePath           string
	PackageName        string
	SourceRelativePath string
	StructSources      StructSourceSlice
}

type StructSourceSlice []*StructSource

type StructSource struct {
	PackageName   string
	TokenPosition token.Position
	TypeSpec      *ast.TypeSpec
	StructType    *ast.StructType
	CommentGroup  *ast.CommentGroup
}

func (s *StructSource) CommentGroupString() string {
	var builder strings.Builder
	for _, comment := range s.CommentGroup.List {
		builder.WriteString(comment.Text + "\n")
	}
	return builder.String()
}

func (s *StructSource) GoString() string {
	return fmt.Sprintf(
		"&StructSource{PackageName: %#v, TokenPosition: %#v, TypeSpec: %#v, StructType: %#v, CommentGroup: %q}",
		s.PackageName,
		s.TokenPosition,
		s.TypeSpec,
		s.StructType,
		s.CommentGroupString(),
	)
}

//nolint:gochecknoglobals
var (
	_GoColumnTagCommentLineRegex     *regexp.Regexp
	_GoColumnTagCommentLineRegexOnce sync.Once
)

const (
	//	                                             _____________ <- 1. comment prefix
	//	                                                             __ <- 2. tag name
	//	                                                                               ___ <- 4. comment suffix
	_GoColumnTagCommentLineRegexFormat       = `^\s*(//+\s*|/\*\s*)?(%s)\s*:\s*(.*)\s*(\*/)?`
	_GoColumnTagCommentLineRegexContentIndex = /*                               ^^ 3. tag value */ 3
)

func GoColumnTagCommentLineRegex(ctx context.Context) *regexp.Regexp {
	cfg := contexts.GenerateConfig(ctx)

	_GoColumnTagCommentLineRegexOnce.Do(func() {
		_GoColumnTagCommentLineRegex = regexp.MustCompile(fmt.Sprintf(_GoColumnTagCommentLineRegexFormat, cfg.GoColumnTag))
	})
	return _GoColumnTagCommentLineRegex
}
