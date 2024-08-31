package testdata_test

import (
	"bytes"
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/hakadoriya/ormgen/internal/lang/go/generator/test/generated/ormgen"
	"github.com/hakadoriya/ormgen/internal/lang/go/generator/test/generated/user"
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
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		queryerContext := newTestQueryerContext()
		orm := user.NewORM()
		buf := bytes.NewBuffer(nil)
		ctx := ormgen.LoggerWithContext(context.Background(), slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})))

		{
			defer func() {
				if r := recover(); r == nil {
					t.Error("❌: panic expected")
				}
				queryerContext.ExecContextArgs = testExecContextArgs{}
			}()
			_, _ = orm.GetUserByPK(ctx, queryerContext, 1)
		}
		{
			_, err := orm.ListUser(ctx, queryerContext, ormgen.Where(ormgen.And(ormgen.Equal("username", "Alice"), ormgen.Equal("city", "Tokyo"), ormgen.Or(ormgen.Equal("group_id", 1), ormgen.Equal("group_id", 2)))))
			if err == nil {
				t.Error("❌: err == nil")
			}
			const expectedQuery = "SELECT user_id, username, address, group_id FROM user WHERE (username = $1 AND city = $2 AND (group_id = $3 OR group_id = $4))"
			actualQuery := queryerContext.QueryContextArgs.query
			if expectedQuery != actualQuery {
				t.Errorf("❌: expected(%s) != actual(%s)", expectedQuery, actualQuery)
			}
		}
	})
}
