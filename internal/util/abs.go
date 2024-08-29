package util

import (
	"fmt"
	"path/filepath"

	"github.com/hakadoriya/ormgen/internal/logs"
)

func Abs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		logs.Stdout.Warn(fmt.Sprintf("failed to get absolute path. use path instead: path=%s: %v", path, err))
		return path
	}
	return abs
}
