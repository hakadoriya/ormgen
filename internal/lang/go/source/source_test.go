package source

import (
	"context"
	"go/ast"
	"testing"
)

func TestStructSource_ExtractTableName(t *testing.T) {
	t.Parallel()

	t.Run("success,table_name", func(t *testing.T) {
		t.Parallel()

		s := &StructSource{CommentGroup: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "// Table is table_name"},
				{Text: "//"},
				{Text: "//db:table table_name"},
			},
		}}

		const expected = "table_name"

		// once
		actual := s.ExtractTableName(context.Background(), goColumnTag)
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}

		// twice (cached)
		actual = s.ExtractTableName(context.Background(), goColumnTag)
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})

	t.Run("error,empty", func(t *testing.T) {
		t.Parallel()

		s := &StructSource{CommentGroup: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "// Table is table_name"},
				{Text: "//"},
				{Text: "//db:table table_name"},
			},
		}}

		const expected = ""
		actual := s.ExtractTableName(context.Background(), "errNoTag")
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})
}
