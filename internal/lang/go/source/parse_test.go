package source

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/ormgen/pkg/apperr"
	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

const (
	goColumnTag     = "db"
	sourcePath      = "testdata"
	sourceFileUser  = "testdata/user/user.go"
	sourceFileOther = "testdata/user/other.go"
	sourceFileGroup = "testdata/group/group.go"
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
		pkgSrcSlice, err := Parse(ctx, []string{sourcePath})
		requirez.NoError(t, err)
		requirez.Equal(t, 2, len(pkgSrcSlice))
		assertz.Equal(t, "group", pkgSrcSlice[0].PackageName)
		requirez.Equal(t, 1, len(pkgSrcSlice[0].FileSources))
		assertz.StringHasSuffix(t, pkgSrcSlice[0].FileSources[0].FilePath, sourceFileGroup)
		assertz.Equal(t, "user", pkgSrcSlice[1].PackageName)
		requirez.Equal(t, 2, len(pkgSrcSlice[1].FileSources))
		assertz.StringHasSuffix(t, pkgSrcSlice[1].FileSources[0].FilePath, sourceFileOther)
		assertz.StringHasSuffix(t, pkgSrcSlice[1].FileSources[1].FilePath, sourceFileUser)
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
		pkgSrcSlice, err := Parse(ctx, []string{sourceFileUser})
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
		pkgSrcSlice, err := Parse(ctx, []string{sourcePath})
		requirez.ErrorIs(t, err, context.Canceled)
		assertz.Nil(t, pkgSrcSlice)
	})
}
