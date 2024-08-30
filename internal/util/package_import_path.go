package util

import (
	"fmt"
	"go/build"
	"path/filepath"
)

func DetectPackageImportPath(path string) (string, error) {
	absDir, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs: path=%s %w", path, err)
	}

	pkg, err := build.ImportDir(absDir, build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("build.ImportDir: path=%s: %w", path, err)
	}

	if pkg.ImportPath == "." {
		// If ImportPath is ".", find the parent package recursively.
		parent, err := DetectPackageImportPath(filepath.Dir(absDir))
		if err != nil {
			return "", fmt.Errorf("DetectPackageImportPath: %w", err)
		}

		return filepath.Join(parent, filepath.Base(absDir)), nil
	}

	return pkg.ImportPath, nil
}
