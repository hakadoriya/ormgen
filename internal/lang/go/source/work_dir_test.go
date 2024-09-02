package source

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"regexp"
	"testing"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/z.go/buildz"
	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

func Test_walkDirFn(t *testing.T) {
	t.Parallel()

	t.Run("error,errArg", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		err := walkDirFn(ctx, "", ".dat", nil)("", nil, io.EOF)
		requirez.ErrorIs(t, err, io.EOF)
	})

	t.Run("error,d.IsDir", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		info, errStat := os.Stat(sourcePath)
		requirez.NoError(t, errStat)
		err := walkDirFn(ctx, "", ".dat", nil)("", fs.FileInfoToDirEntry(info), nil)
		requirez.NoError(t, err)
	})

	t.Run("error,no-such-file-or-directory.dat", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		trace := bytes.NewBuffer(nil)
		ctx = contexts.WithTrace(ctx, logs.NewTrace(trace, slog.LevelDebug))
		info, errStat := os.Stat(sourceFileUser)
		requirez.NoError(t, errStat)
		err := walkDirFn(ctx, "", ".dat", nil)("no-such-file-or-directory.dat", fs.FileInfoToDirEntry(info), nil)
		requirez.NoError(t, err)
		actualLog := trace.String()
		assertz.StringMatchRegex(t, actualLog, regexp.MustCompile(`"severity":"DEBUG","caller":"source/work_dir.go:\d+","msg":"skip: path=no-such-file-or-directory.dat: parser.ParseFile=open no-such-file-or-directory.dat: no such file or directory: no source found","app":"ormgen"}`))
	})

	t.Run("error,buildz.ErrReachedRootDirectory", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		trace := bytes.NewBuffer(nil)
		ctx = contexts.WithTrace(ctx, logs.NewTrace(trace, slog.LevelDebug))
		info, errStat := os.Stat(sourceFileUser)
		requirez.NoError(t, errStat)
		err := walkDirFn(ctx, "/", fileExt, nil)(sourceFileUser, fs.FileInfoToDirEntry(info), nil)
		requirez.ErrorIs(t, err, buildz.ErrReachedRootDirectory)
	})

	t.Run("error,filepath.Rel", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithGenerateConfig(context.Background(), &config.GenerateConfig{
			GoColumnTag: goColumnTag,
		})
		trace := bytes.NewBuffer(nil)
		ctx = contexts.WithTrace(ctx, logs.NewTrace(trace, slog.LevelDebug))
		info, errStat := os.Stat(sourceFileUser)
		requirez.NoError(t, errStat)
		err := walkDirFn(ctx, "..", fileExt, nil)(sourceFileUser, fs.FileInfoToDirEntry(info), nil)
		requirez.ErrorContains(t, err, `Rel: can't make testdata/user/user.go relative to ..`)
	})

}