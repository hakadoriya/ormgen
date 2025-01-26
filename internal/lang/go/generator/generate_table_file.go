package generator

import (
	"context"
	"os"
	"strings"

	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/otelz/tracez"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/ormgen/internal/logs"
)

func generateTableFile(ctx context.Context, fileSource *source.FileSource) (err error) {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	cfg := contexts.GenerateConfig(ctx)

	// table file
	dbGenFilePath := strings.TrimSuffix(fileSource.FilePath, ".go") + "." + cfg.GoColumnTag + ".gen.go"
	logs.Stdout.Debug("create file: file=" + dbGenFilePath)
	var dbGenFile *os.File
	if err := tracez.StartFuncWithSpanNameSuffix(ctx, "os.Create", func(_ context.Context) (err error) {
		dbGenFile, err = os.Create(dbGenFilePath)
		return err //nolint:wrapcheck
	}); err != nil {
		return errorz.Errorf("os.Create: %w", err)
	}
	defer dbGenFile.Close()

	if err := fprintTableMethods(ctx, dbGenFile, fileSource); err != nil {
		return errorz.Errorf("fprintTableMethods: %w", err)
	}

	return nil
}
