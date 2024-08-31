package sandbox

import (
	"strconv"
	"testing"
)

func TestQueryOption(t *testing.T) {
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

	QueryPrefix("@{HOGE=hoge}").apply(c)
	Where(
		And(
			Equal("name", "hoge"),
			LessThan("age", 20),
			Or(
				NotIn("group_id", 1, 2, 3),
				Equal("is_admin", true),
			),
		),
	).apply(c)
	OrderByDesc("created_at").apply(c)
	Limit(10).apply(c)

	query, args := c.ToSQL("SELECT * FROM user WHERE group_id = $1")

	t.Logf("\n" + query)
	t.Logf("\n%v", args)
}
