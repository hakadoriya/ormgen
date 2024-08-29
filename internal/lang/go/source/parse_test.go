package source

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/ormgen/pkg/apperr"
)

const (
	goColumnTag           = "db"
	validSourcePath       = "testdata/valid"
	validSourceFileUser   = "testdata/valid/user/user.go"
	validSourceFileOther  = "testdata/valid/user/other.go"
	validSourceFileGroup  = "testdata/valid/group/group.go"
	invalidSourcePath     = "testdata/invalid"
	invalidSourceFileUser = "testdata/invalid/user/user.go"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("success,Parse", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		stderr := bytes.NewBuffer(nil)
		ctx = contexts.WithStderr(ctx, logs.NewStderr(stderr, slog.LevelDebug))
		pkgSrcSlice, err := Parse(ctx, []string{validSourcePath})
		requirez.NoError(t, err)
		requirez.Equal(t, 2, len(pkgSrcSlice))
		assertz.Equal(t, "group", pkgSrcSlice[0].PackageName)
		requirez.Equal(t, 1, len(pkgSrcSlice[0].FileSources))
		assertz.StringHasSuffix(t, pkgSrcSlice[0].FileSources[0].FilePath, validSourceFileGroup)
		assertz.Equal(t, "user", pkgSrcSlice[1].PackageName)
		requirez.Equal(t, 2, len(pkgSrcSlice[1].FileSources))
		assertz.StringHasSuffix(t, pkgSrcSlice[1].FileSources[0].FilePath, validSourceFileOther)
		assertz.StringHasSuffix(t, pkgSrcSlice[1].FileSources[1].FilePath, validSourceFileUser)
	})

	t.Run("error,apperr.ErrInvalidArguments", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		pkgSrcSlice, err := Parse(ctx, nil)
		requirez.ErrorIs(t, err, apperr.ErrInvalidArguments)
		assertz.ErrorContains(t, err, `invalid number of arguments; expected 1, got 0:`)
		assertz.Nil(t, pkgSrcSlice)
	})

	t.Run("error,os.Stat", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		pkgSrcSlice, err := Parse(ctx, []string{"no-such-file-or-directory"})
		requirez.ErrorContains(t, err, `: no such file or directory`)
		assertz.Nil(t, pkgSrcSlice)
	})

	t.Run("error,apperr.ErrSourcePathIsNotDirectory", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		pkgSrcSlice, err := Parse(ctx, []string{validSourceFileUser})
		requirez.ErrorIs(t, err, apperr.ErrSourcePathIsNotDirectory)
		assertz.Nil(t, pkgSrcSlice)
	})

	t.Run("error,context.Canceled", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		pkgSrcSlice, err := Parse(ctx, []string{validSourcePath})
		requirez.ErrorIs(t, err, context.Canceled)
		assertz.Nil(t, pkgSrcSlice)
	})
}
