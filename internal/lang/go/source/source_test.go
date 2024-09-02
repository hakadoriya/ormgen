package source

import (
	"context"
	"go/ast"
	"testing"
)

func TestStructSource_ExtractTableName(t *testing.T) {
	t.Parallel()

	t.Run("error,", func(t *testing.T) {
		t.Parallel()

		s := &StructSource{CommentGroup: &ast.CommentGroup{}}

		const expected = ""
		actual := s.ExtractTableName(context.Background(), goColumnTag)
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})
}
