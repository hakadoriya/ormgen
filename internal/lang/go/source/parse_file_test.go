package source

import (
	"bytes"
	"context"
	"log/slog"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/ormgen/pkg/apperr"
)

func Test_parseFile(t *testing.T) {
	t.Parallel()

	t.Run("success,parseFile,testdata", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		stdout := bytes.NewBuffer(nil)
		ctx = contexts.WithStdout(ctx, logs.NewStdout(stdout, slog.LevelDebug))
		stderr := bytes.NewBuffer(nil)
		ctx = contexts.WithStderr(ctx, logs.NewStderr(stderr, slog.LevelDebug))

		fileSource, err := parseFile(ctx, validSourcePath, validSourceFileUser)
		requirez.NoError(t, err)
		assertz.Equal(t, validSourceFileUser, fileSource.FilePath)
		assertz.Equal(t, "user", fileSource.PackageName)
		assertz.Equal(t, 2, len(fileSource.StructSources))
		assertz.Equal(t, "User", fileSource.StructSources[0].TypeSpec.Name.Name)
		assertz.Equal(t, "user", fileSource.StructSources[0].ExtractTableName(ctx, goColumnTag))
		assertz.Equal(t, "AdminUser", fileSource.StructSources[1].TypeSpec.Name.Name)
		assertz.Equal(t, "admin_user", fileSource.StructSources[1].ExtractTableName(ctx, goColumnTag))
		assertz.StringMatchRegex(t, stdout.String(), regexp.MustCompile(`"severity":"DEBUG","caller":"source/parse_file.go:\d+","msg":"parse file: file=testdata/valid/user/user.go","logName":"stdout","app":"ormgen"}`))
		assertz.StringMatchRegex(t, stdout.String(), regexp.MustCompile(`"severity":"DEBUG","caller":"source/parse_file.go:\d+","msg":"found struct source: file=testdata/valid/user/user.go:\d+:\d+, tag=db:table, type=User","logName":"stdout","app":"ormgen"}`))
		assertz.StringMatchRegex(t, stdout.String(), regexp.MustCompile(`"severity":"DEBUG","caller":"source/parse_file.go:\d+","msg":"found struct source: file=testdata/valid/user/user.go:\d+:\d+, tag=db:table, type=AdminUser","logName":"stdout","app":"ormgen"}`))
	})

	t.Run("error,apperr.ErrNoSourceFound", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		fileSource, err := parseFile(ctx, validSourcePath, "no-such-file-or-directory")
		requirez.ErrorIs(t, err, apperr.ErrNoSourceFound)
		assertz.Nil(t, fileSource)
	})

	t.Run("error,filepath.Rel", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		fileSource, err := parseFile(ctx, string(filepath.Separator), validSourceFileUser)
		requirez.ErrorContains(t, err, `filepath.Rel: Rel: can't make testdata/valid/user/user.go relative to /`)
		assertz.Nil(t, fileSource)
	})

	t.Run("error,apperr.ErrInvalidAnnotation", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		fileSource, err := parseFile(ctx, invalidSourcePath, invalidSourceFileUser)
		requirez.ErrorIs(t, err, apperr.ErrInvalidAnnotation)
		assertz.Nil(t, fileSource)
	})

	t.Run("error,apperr.ErrNoSourceFound", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: "no-such-tag",
		})
		fileSource, err := parseFile(ctx, validSourcePath, validSourceFileUser)
		requirez.ErrorIs(t, err, apperr.ErrNoSourceFound)
		assertz.Nil(t, fileSource)
	})
}
