package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/hakadoriya/ormgen/internal/apperr"
	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/entrypoint"

	"github.com/hakadoriya/z.go/buildinfoz"
	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/logz/slogz"
)

func main() {
	slog.SetDefault(slog.New(slogz.NewHandler(os.Stdout, slog.LevelDebug)))

	exitCode, err := exec()
	if err != nil {
		apperr.Log.Error(fmt.Sprintf("exit %d", exitCode), slog.Any("error", err))
	}
	os.Exit(exitCode)
}

func exec() (exitCode int, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*
		_ = []cliz.Option{
			// Common
			&cliz.BoolOption{Name: config.OptionDebug, Env: config.EnvDebug, Description: "Enable debug mode"},
			&cliz.StringOption{Name: config.OptionDialect, Env: config.EnvKeyDialect, Default: "dialect", Description: "dialect for DML"},
			&cliz.StringOption{Name: config.OptionLanguage, Env: config.EnvKeyLanguage, Default: "go", Description: "programming language to generate ORM"},
			// Go
			&cliz.StringOption{Name: config.OptionGoColumnTag, Env: config.EnvGoColumnTag, Default: "db", Description: "column annotation key for Go struct tag"},
			&cliz.StringOption{Name: config.OptionGoPKTag, Env: config.EnvKeyGoPKTag, Default: "pk", Description: "primary key annotation key for Go struct tag"},
			&cliz.StringOption{Name: config.OptionGoHasOneTag, Env: config.EnvKeyGoHasOneTag, Default: "hasOne", Description: "\"hasOne\" annotation key for Go struct tag"},
			&cliz.StringOption{Name: config.OptionGoHasManyTag, Env: config.EnvKeyGoHasManyTag, Default: "hasMany", Description: "\"hasMany\" annotation key for Go struct tag"},
			&cliz.StringOption{Name: config.OptionGoTableNameMethod, Env: config.EnvKeyGoTableNameMethod, Default: "TableName", Description: "method name for table"},
			&cliz.StringOption{Name: config.OptionGoColumnNameMethodPrefix, Env: config.EnvKeyGoColumnNameMethodPrefix, Default: "ColumnName_", Description: "method name for columns"},
			&cliz.StringOption{Name: config.OptionGoColumnsNameMethod, Env: config.EnvKeyGoColumnsNameMethod, Default: "ColumnsName", Description: "method prefix for column"},
			&cliz.StringOption{Name: config.OptionGoSliceTypeSuffix, Env: config.EnvKeyGoSliceTypeSuffix, Default: "Slice", Description: "suffix for slice type"},
			&cliz.StringOption{Name: config.OptionGoORMOutputPath, Env: config.EnvKeyGoORMOutputPath, Default: "ormgen", Description: "output path of ORM."},
			&cliz.StringOption{Name: config.OptionGoORMPackageName, Env: config.EnvKeyGoORMPackageName, Default: "", Description: "package name for ORM. If empty, use the base name of the output path."},
			&cliz.StringOption{Name: config.OptionGoORMStructPackageImportPath, Env: config.EnvKeyGoORMStructPackageImportPath, Description: "package import path of ORM target struct. If empty, try to detect automatically."},
			&cliz.StringOption{Name: config.OptionGoORMInterfaceName, Env: config.EnvKeyGoORMInterfaceName, Default: "ORM", Description: "interface type name for ORM"},
			&cliz.StringOption{Name: config.OptionGoORMStructName, Env: config.EnvKeyGoORMStructName, Default: "_ORM", Description: "struct name for ORM"},
		}
	*/

	generateOpts, err := cliz.MarshalOptions(new(config.Generate))
	if err != nil {
		return 1, errorz.Errorf("cliz.MarshalOptions: %w", err)
	}

	//nolint:exhaustruct
	cmd := &cliz.Command{
		Name: config.AppName,
		SubCommands: []*cliz.Command{
			{
				Name:            "generate",
				Usage:           "generate <source>",
				Description:     "Generate ORM from annotated code",
				Options:         generateOpts,
				PreHookExecFunc: config.Load,
				ExecFunc:        entrypoint.Generate,
			},
			{
				Name:        "version",
				Description: "Show version information",
				ExecFunc: func(c *cliz.Command, _ []string) error {
					return buildinfoz.Fprint(c.Stdout())
				},
			},
		},
	}

	if err := cmd.Exec(ctx, os.Args); err != nil && !cliz.IsHelp(err) {
		return 1, errorz.Errorf("cmd.Exec: %w", err)
	}

	return 0, nil
}
