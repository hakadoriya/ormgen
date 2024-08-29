package iface

import (
	"errors"
	"os"
	"testing"
)

func TestPkg_FilepathAbs(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		p := &Pkg{
			FilepathAbsFunc: func(path string) (string, error) {
				return "", os.ErrNotExist
			},
		}

		_, err := p.FilepathAbs("path")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("❌: expected(%v) != actual(%v)", os.ErrNotExist, err)
		}
	})
}
