package ormgen

import (
	"reflect"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
)

func TestQueryOption(t *testing.T) {
	const query = "SELECT * FROM user WHERE group_id = $1" // placeholder 1

	c := new(QueryConfig)

	opts := []QueryOption{
		WithPlaceholderGenerator(PlaceholderGeneratorMap["postgres"]),
		QueryPrefix("@{HOGE=hoge}"),
		Where(
			And(
				Equal("name", "Alice"), // placeholder 2
				LessThan("age", 20),    // placeholder 3
				Or(
					NotIn("group_id", 1, 2, 3), // placeholder 4, 5, 6
					Equal("is_admin", true),    // placeholder 7
				),
			),
		),
		OrderByDesc("created_at"),
		Limit(10), // placeholder 8
	}

	for _, o := range opts {
		o.ApplyQueryOption(c)
	}

	q, args := c.ToSQL(query, 2)
	t.Logf("\n" + q)
	t.Logf("\n%#v", args)

	const expectedQuery = `@{HOGE=hoge} SELECT * FROM user WHERE group_id = $1 WHERE (name = $2 AND age < $3 AND (group_id NOT IN ($4, $5, $6) OR is_admin = $7)) ORDER BY created_at DESC LIMIT $8`
	actualQuery := q
	assertz.Equal(t, expectedQuery, actualQuery)

	expectedArgs := []interface{}{"Alice", 20, 1, 2, 3, true, 10}
	actualArgs := args
	if !reflect.DeepEqual(expectedArgs, actualArgs) {
		t.Errorf("❌: expected(%#v) != actual(%#v)", expectedArgs, actualArgs)
	}

}
