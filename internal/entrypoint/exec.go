package entrypoint

import (
	"context"
	"os"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/errorz"
)

func Exec(ctx context.Context) (exitCode int, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generateOpts, err := cliz.MarshalOptions(new(config.GenerateConfig))
	if err != nil {
		return 1, errorz.Errorf("cliz.MarshalOptions: %w", err)
	}

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
				ExecFunc:        Generate,
			},
			{
				Name:        "version",
				Description: "Show version information",
				ExecFunc:    Version,
			},
		},
	}

	if err := cmd.Exec(ctx, os.Args); err != nil && !cliz.IsHelp(err) {
		return 1, errorz.Errorf("cmd.Exec: %w", err)
	}

	return 0, nil
}
