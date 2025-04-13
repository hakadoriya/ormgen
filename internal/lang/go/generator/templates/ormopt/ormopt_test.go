package ormopt

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
)

// func TestPlaceholderGeneratorMap(t *testing.T) {
//	t.Parallel()
//
//	t.Run("success,postgres", func(t *testing.T) {
//		t.Parallel()
//
//		actual := PlaceholderGeneratorMap["postgres"](10)
//		requirez.Equal(t, "$10", actual)
//	})
//
//	t.Run("success,cockroach", func(t *testing.T) {
//		t.Parallel()
//
//		actual := PlaceholderGeneratorMap["cockroach"](10)
//		requirez.Equal(t, "$10", actual)
//	})
//
//	t.Run("success,mysql", func(t *testing.T) {
//		t.Parallel()
//
//		actual := PlaceholderGeneratorMap["mysql"](10)
//		requirez.Equal(t, "?", actual)
//	})
//
//	t.Run("success,sqlite3", func(t *testing.T) {
//		t.Parallel()
//
//		actual := PlaceholderGeneratorMap["sqlite3"](10)
//		requirez.Equal(t, "?", actual)
//	})
//
//	t.Run("success,spanner", func(t *testing.T) {
//		t.Parallel()
//
//		actual := PlaceholderGeneratorMap["spanner"](10)
//		requirez.Equal(t, "?", actual)
//	})
//
//	t.Run("error,empty", func(t *testing.T) {
//		t.Parallel()
//
//		defer func() {
//			if r := recover(); r == nil {
//				t.Errorf("❌: expected panic")
//			}
//		}()
//
//		actual := PlaceholderGeneratorMap[""](10)
//		requirez.Equal(t, "", actual)
//	})
//}

func TestQueryOption(t *testing.T) {
	t.Parallel()

	t.Run("success,Case1", func(t *testing.T) {
		t.Parallel()

		const query = "SELECT * FROM user WHERE group_id = $1" // placeholder 1

		c := new(QueryConfig)

		opts := []QueryOption{
			WithPlaceholderGenerator(func(i int) string { return "$" + strconv.Itoa(i) }),
			QueryPrefix("@{HOGE=hoge}"),
			Where(
				And(
					Equal("name", "Alice"), // placeholder 2
					LessThan("age", 20),    // placeholder 3
					Or(
						In("group_id", 1, 2, 3),       // placeholder 4, 5, 6
						NotIn("group_id", 4, 5),       // placeholder 7, 8
						NotEqual("is_admin", true),    // placeholder 9
						GreaterThan("age", 9),         // placeholder 10
						GreaterThanOrEqual("age", 10), // placeholder 11
						LessThanOrEqual("age", 39),    // placeholder 12
						LessThan("age", 40),           // placeholder 13
					),
					Like("name", "A%"),    // placeholder 14
					NotLike("name", "B%"), // placeholder 15
					IsNull("address"),
					IsNotNull("email"),
					Between("created_at", "2021-01-01", "2021-12-31"),    // placeholder 16, 17
					NotBetween("updated_at", "2022-01-01", "2022-12-31"), // placeholder 18, 19
				),
			),
			OrderByDesc("created_at"),
			Limit(10), // placeholder 20
		}

		for _, o := range opts {
			o.ApplyQueryOption(c)
		}

		q, args := c.ToSQL(query, 2)
		t.Logf("\n%s", q)
		t.Logf("\n%#v", args)

		const expectedQuery = `@{HOGE=hoge} SELECT * FROM user WHERE group_id = $1 WHERE (name = $2 AND age < $3 AND (group_id IN ($4, $5, $6) OR group_id NOT IN ($7, $8) OR is_admin <> $9 OR age > $10 OR age >= $11 OR age <= $12 OR age < $13) AND name LIKE $14 AND name NOT LIKE $15 AND address IS NULL AND email IS NOT NULL AND created_at BETWEEN $16 AND $17 AND updated_at NOT BETWEEN $18 AND $19) ORDER BY created_at DESC LIMIT $20`
		actualQuery := q
		assertz.Equal(t, expectedQuery, actualQuery)

		expectedArgs := []interface{}{"Alice", 20, 1, 2, 3, 4, 5, true, 9, 10, 39, 40, "A%", "B%", "2021-01-01", "2021-12-31", "2022-01-01", "2022-12-31", 10}
		actualArgs := args
		if !reflect.DeepEqual(expectedArgs, actualArgs) {
			t.Errorf("❌: expected(%#v) != actual(%#v)", expectedArgs, actualArgs)
		}
	})

	t.Run("success,Case1", func(t *testing.T) {
		t.Parallel()

		const query = "SELECT * FROM user WHERE group_id = $1" // placeholder 1

		c := new(QueryConfig)

		opts := []QueryOption{
			OrderBy("created_at"),
		}

		for _, o := range opts {
			o.ApplyQueryOption(c)
		}

		q, args := c.ToSQL(query, 2)
		t.Logf("query:\n%s", q)
		t.Logf("args:\n%#v", args)

		const expectedQuery = `SELECT * FROM user WHERE group_id = $1 ORDER BY created_at`
		actualQuery := q
		assertz.Equal(t, expectedQuery, actualQuery)

		expectedArgs := []interface{}(nil)
		actualArgs := args
		if !reflect.DeepEqual(expectedArgs, actualArgs) {
			t.Errorf("❌: expected(%#v) != actual(%#v)", expectedArgs, actualArgs)
		}
	})
}
