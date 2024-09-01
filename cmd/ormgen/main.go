package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/ormgen/pkg/ormgen"
)

func main() {
	exitCode, err := ormgen.Exec(context.Background(), os.Args)
	if err != nil {
		logs.Stderr.Error(fmt.Sprintf("exit %d", exitCode), slog.Any("error", err))
	}
	os.Exit(exitCode)
}
