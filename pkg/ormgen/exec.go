package ormgen

import (
	"context"

	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/mustz"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/entrypoint"
)

func Exec(ctx context.Context, osArgs []string) (exitCode int, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	generateOpts := mustz.One(cliz.MarshalOptions(new(config.GenerateConfig)))

	//nolint:exhaustruct
	cmd := &cliz.Command{
		Name: consts.AppName,
		SubCommands: []*cliz.Command{
			{
				Name:            "generate",
				Usage:           "generate <SOURCE DIRECTORY>",
				Description:     "Generate ORM from annotated code",
				Options:         generateOpts,
				PreHookExecFunc: config.GeneratePreHookExec,
				ExecFunc:        entrypoint.Generate,
			},
			{
				Name:        "version",
				Description: "Show version information",
				ExecFunc:    entrypoint.Version,
			},
		},
	}

	if err := cmd.Exec(ctx, osArgs); err != nil && !cliz.IsHelp(err) {
		return 1, errorz.Errorf("cmd.Exec: %w", err)
	}

	return 0, nil
}
