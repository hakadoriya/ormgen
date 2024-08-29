package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hakadoriya/z.go/logz/slogz"

	"github.com/hakadoriya/ormgen/internal/logs"
	"github.com/hakadoriya/ormgen/pkg/ormgen"
)

func main() {
	exitCode, err := ormgen.Exec(context.Background(), os.Args)
	if err != nil {
		logs.Stderr.Error(fmt.Sprintf("exit %d", exitCode), slogz.Error(err))
	}
	os.Exit(exitCode)
}
