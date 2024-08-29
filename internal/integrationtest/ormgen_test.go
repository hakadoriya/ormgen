package ormgen_test

import (
	"bytes"
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/hakadoriya/ormgen/examples/generated/postgres/ormopt"
	"github.com/hakadoriya/ormgen/examples/generated/postgres/user"
	model_user "github.com/hakadoriya/ormgen/examples/model/user"
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

func TestListUser(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		queryerContext := newTestQueryerContext()
		orm := user.NewORM()
		buf := bytes.NewBuffer(nil)
		ctx := ormopt.LoggerWithContext(context.Background(), slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})))

		{
			_, err := orm.ListUser(ctx, queryerContext,
				ormopt.Where(
					ormopt.And(
						ormopt.Equal("username", "Alice"), // placeholder $1
						ormopt.NotEqual("city", "Tokyo"),  // placeholder $2
						ormopt.Or(
							ormopt.In("group_id", 1, 2, 3),     // placeholder $3 $4 $5
							ormopt.NotIn("attribute_id", 4, 5), // placeholder $6 $7
						),
					),
				),
				ormopt.OrderByDesc("created_at"),
				ormopt.Limit(10), // placeholder $8
			)
			if err == nil {
				t.Error("❌: err == nil")
			}
			actualQuery := queryerContext.QueryContextArgs.query
			t.Logf("📝: query:\n%s", actualQuery)
			const expectedQuery = `SELECT user_id, username, address, group_id FROM user WHERE (username = $1 AND city <> $2 AND (group_id IN ($3, $4, $5) OR attribute_id NOT IN ($6, $7))) ORDER BY created_at DESC LIMIT $8`
			if expectedQuery != actualQuery {
				t.Errorf("❌: expected(%s) != actual(%s)", expectedQuery, actualQuery)
			}
		}
	})
}

func TestBulkInsertUser(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		queryerContext := newTestQueryerContext()
		orm := user.NewORM()
		buf := bytes.NewBuffer(nil)
		ctx := ormopt.LoggerWithContext(context.Background(), slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})))

		user.DefaultBulkInsertMaxPlaceholdersPerQuery = 50
		users := func() []*model_user.User {
			users := make([]*model_user.User, 0)
			for i := 1; i <= 25; i++ {
				users = append(users, &model_user.User{
					UserID:   i,
					Username: "Alice",
					Address:  "Tokyo",
					GroupID:  1,
				})
			}
			return users
		}()

		{
			err := orm.BulkInsertUser(ctx, queryerContext, users)
			if err == nil {
				t.Error("❌: err == nil")
			}
			actualQuery := queryerContext.ExecContextArgs.query
			t.Logf("📝: query:\n%s", actualQuery)
			const expectedQuery = `INSERT INTO user (user_id, username, address, group_id) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ($9, $10, $11, $12), ($13, $14, $15, $16), ($17, $18, $19, $20), ($21, $22, $23, $24), ($25, $26, $27, $28), ($29, $30, $31, $32), ($33, $34, $35, $36), ($37, $38, $39, $40), ($41, $42, $43, $44), ($45, $46, $47, $48)`
			if expectedQuery != actualQuery {
				t.Errorf("❌: expected(%s) != actual(%s)", expectedQuery, actualQuery)
			}
		}
	})
}
