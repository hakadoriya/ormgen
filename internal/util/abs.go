package util

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hakadoriya/ormgen/internal/contexts"
)

func Abs(ctx context.Context, path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		contexts.Stdout(ctx).Warn(fmt.Sprintf("failed to get absolute path. use path instead: path=%s: %v", path, err))
		return path
	}
	return abs
}
