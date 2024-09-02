package ormgen_test

import (
	"bytes"
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/hakadoriya/ormgen/example/generated/postgres/ormgen"
	"github.com/hakadoriya/ormgen/example/generated/postgres/user"
)

type (
	testQueryerContext struct {
		ExecContextFunc     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		ExecContextArgs     testExecContextArgs
		PrepareContextFunc  func(ctx context.Context, query string) (*sql.Stmt, error)
		PrepareContextArgs  testPrepareContextArgs
		QueryContextFunc    func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryContextArgs    testQueryContextArgs
		QueryRowContextFunc func(ctx context.Context, query string, args ...interface{}) *sql.Row
		QueryRowContextArgs testQueryRowContextArgs
	}
	testExecContextArgs struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	testPrepareContextArgs struct {
		ctx   context.Context
		query string
	}
	testQueryContextArgs struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	testQueryRowContextArgs struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
)

func (s *testQueryerContext) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	s.ExecContextArgs = testExecContextArgs{ctx, query, args}
	return s.ExecContextFunc(ctx, query, args...)
}

func (s *testQueryerContext) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	s.PrepareContextArgs = testPrepareContextArgs{ctx, query}
	return s.PrepareContextFunc(ctx, query)
}

func (s *testQueryerContext) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	s.QueryContextArgs = testQueryContextArgs{ctx, query, args}
	return s.QueryContextFunc(ctx, query, args...)
}

func (s *testQueryerContext) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	s.QueryRowContextArgs = testQueryRowContextArgs{ctx, query, args}
	return s.QueryRowContextFunc(ctx, query, args...)
}

func newTestQueryerContext() *testQueryerContext {
	return &testQueryerContext{
		ExecContextFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return nil, sql.ErrTxDone
		},
		PrepareContextFunc: func(ctx context.Context, query string) (*sql.Stmt, error) {
			return nil, sql.ErrTxDone
		},
		QueryRowContextFunc: func(ctx context.Context, query string, args ...interface{}) *sql.Row {
			return &sql.Row{}
		},
		QueryContextFunc: func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
			return nil, sql.ErrTxDone
		},
	}
}

func TestOrmGen(t *testing.T) {
	t.Run("success,", func(t *testing.T) {
		queryerContext := newTestQueryerContext()
		orm := user.NewORM()
		buf := bytes.NewBuffer(nil)
		ctx := ormgen.LoggerWithContext(context.Background(), slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})))

		{
			_, err := orm.ListUser(ctx, queryerContext,
				ormgen.Where(
					ormgen.And(
						ormgen.Equal("username", "Alice"), // placeholder $1
						ormgen.NotEqual("city", "Tokyo"),  // placeholder $2
						ormgen.Or(
							ormgen.In("group_id", 1, 2, 3),     // placeholder $3 $4 $5
							ormgen.NotIn("attribute_id", 4, 5), // placeholder $6 $7
						),
					),
				),
				ormgen.OrderByDesc("created_at"),
				ormgen.Limit(10), // placeholder $8
			)
			if err == nil {
				t.Error("❌: err == nil")
			}
			t.Logf("📝: query:\n%s", queryerContext.QueryContextArgs.query)
			const expectedQuery = `SELECT user_id, username, address, group_id FROM user WHERE (username = $1 AND city <> $2 AND (group_id IN ($3, $4, $5) OR attribute_id NOT IN ($6, $7))) ORDER BY created_at DESC LIMIT $8`
			actualQuery := queryerContext.QueryContextArgs.query
			if expectedQuery != actualQuery {
				t.Errorf("❌: expected(%s) != actual(%s)", expectedQuery, actualQuery)
			}
		}
	})
}
