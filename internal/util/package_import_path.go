package util

import (
	"fmt"
	"go/build"
	"path/filepath"

	"github.com/hakadoriya/ormgen/pkg/apperr"
)

func DetectPackageImportPath(path string) (string, error) {
	if path == "" || path == "/" {
		return "", fmt.Errorf("path = %q: %w", path, apperr.ErrEmpty)
	}

	absDir, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs: path=%s: %w", path, err)
	}

	pkg, err := build.ImportDir(absDir, build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("build.ImportDir: path=%s: %w", path, err)
	}

	if pkg.ImportPath == "." {
		// If ImportPath is ".", find the parent package recursively.
		basename := filepath.Base(absDir)
		parentDir := filepath.Dir(absDir)
		parent, err := DetectPackageImportPath(parentDir)
		if err != nil {
			return "", fmt.Errorf("DetectPackageImportPath: %w", err)
		}

		return filepath.Join(parent, basename), nil
	}

	return pkg.ImportPath, nil
}
