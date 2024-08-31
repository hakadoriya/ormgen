package templates

import (
	"strconv"
	"testing"
)

func TestQueryOption(t *testing.T) {
	const query = "SELECT * FROM user WHERE group_id = $1" // placeholder 1

	c := &queryConfig{
		PlaceholderGenGen: func(placeholderStartAt *int, queryArgs *[]interface{}) PlaceholderGen {
			return func(args ...interface{}) string {
				defer func() {
					*placeholderStartAt++
					*queryArgs = append(*queryArgs, args...)
				}()
				return "$" + strconv.Itoa(*placeholderStartAt)
			}
		},
		PlaceholderStartAt: 2,
	}

	opts := []QueryOption{
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
		o.applyQueryOption(c)
	}

	q, args := c.ToSQL(query)

	t.Logf("\n" + q)

	t.Logf("\n%#v", args)
}
