package util

import (
	"fmt"
	"path/filepath"

	"log/slog"
)

func Abs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		slog.Default().Warn(fmt.Sprintf("failed to get absolute path. use path instead: path=%s: %v", path, err))
		return path
	}
	return abs
}
