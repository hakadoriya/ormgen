package source

import (
	"context"
	"io"
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func Test_walkDirFn(t *testing.T) {
	t.Parallel()

	t.Run("error,", func(t *testing.T) {
		t.Parallel()

		err := walkDirFn(context.Background(), "", nil)("", nil, io.EOF)
		requirez.ErrorIs(t, err, io.EOF)
	})
}
