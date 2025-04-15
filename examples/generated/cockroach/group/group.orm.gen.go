// Code generated by ormgen; DO NOT EDIT.
//
// source: group/group.go
package group

import (
	"context"
	"fmt"
	"strings"

	ormcommon "github.com/hakadoriya/ormgen/examples/generated/cockroach/ormcommon"
	ormopt "github.com/hakadoriya/ormgen/examples/generated/cockroach/ormopt"

	group_ "github.com/hakadoriya/ormgen/examples/model/group"
)

const InsertGroupQuery = `INSERT INTO group (id, name) VALUES ($1, $2)`

func (s *_ORM) InsertGroup(ctx context.Context, queryerContext ormcommon.QueryerContext, group *group_.Group) error {
	ormcommon.LoggerFromContext(ctx).Debug(InsertGroupQuery)
	_, err := queryerContext.ExecContext(ctx, InsertGroupQuery, group.ID, group.Name)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const BulkInsertGroupQueryPrefix = `INSERT INTO group (id, name) VALUES `

func (s *_ORM) BulkInsertGroup(ctx context.Context, queryerContext ormcommon.QueryerContext, groupSlice []*group_.Group) error {
	if len(groupSlice) == 0 {
		return nil
	}

	// Calculate the number of placeholders per row and the maximum number of rows per query
	const (
		placeholderStartAt = 1
		placeholdersPerRow = 2
	)

	maxRowsPerQuery := DefaultBulkInsertMaxPlaceholdersPerQuery / placeholdersPerRow
	placeholderIdx := placeholderStartAt
	for i := 0; i < len(groupSlice); i += maxRowsPerQuery {
		end := i + maxRowsPerQuery
		if end > len(groupSlice) {
			end = len(groupSlice)
		}

		chunk := groupSlice[i:end]
		placeholders := make([]string, len(chunk))
		args := make([]interface{}, 0, len(chunk)*placeholdersPerRow)

		for j := range chunk {
			placeholders[j] += "("
			for i := range placeholdersPerRow {
				if i > 0 {
					placeholders[j] += ", "
				}
				placeholders[j] += DefaultPlaceholderGenerator(placeholderIdx)
				placeholderIdx++
			}
			placeholders[j] += ")"
			args = append(args, chunk[j].ID, chunk[j].Name)
		}

		query := BulkInsertGroupQueryPrefix + strings.Join(placeholders, ", ")
		ormcommon.LoggerFromContext(ctx).Debug(query)

		_, err := queryerContext.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
		}
	}

	return nil
}

const GetGroupByPKQuery = `SELECT id, name FROM group WHERE id = $1`

func (s *_ORM) GetGroupByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, id int) (*group_.Group, error) {
	ormcommon.LoggerFromContext(ctx).Debug(GetGroupByPKQuery)
	row := queryerContext.QueryRowContext(ctx, GetGroupByPKQuery, id)
	group := new(group_.Group)
	err := row.Scan(&group.ID, &group.Name)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return group, nil
}

const LockGroupByPKQuery = `SELECT id, name FROM group WHERE id = $1 FOR UPDATE`

func (s *_ORM) LockGroupByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, id int) (*group_.Group, error) {
	ormcommon.LoggerFromContext(ctx).Debug(LockGroupByPKQuery)
	row := queryerContext.QueryRowContext(ctx, LockGroupByPKQuery, id)
	group := new(group_.Group)
	err := row.Scan(&group.ID, &group.Name)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return group, nil
}

const ListGroupQuery = `SELECT id, name FROM group`

func (s *_ORM) ListGroup(ctx context.Context, queryerContext ormcommon.QueryerContext, opts ...ormopt.QueryOption) (group_.GroupSlice, error) {
	config := new(ormopt.QueryConfig)
	ormopt.WithPlaceholderGenerator(DefaultPlaceholderGenerator).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyQueryOption(config)
	}
	query, args := config.ToSQL(ListGroupQuery, 1)
	ormcommon.LoggerFromContext(ctx).Debug(query)
	rows, err := queryerContext.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("queryerContext.QueryContext: %w", s.HandleError(ctx, err))
	}
	var groupSlice group_.GroupSlice
	for rows.Next() {
		group := new(group_.Group)
		err := rows.Scan(&group.ID, &group.Name)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", s.HandleError(ctx, err))
		}
		groupSlice = append(groupSlice, group)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("rows.Close: %w", s.HandleError(ctx, err))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", s.HandleError(ctx, err))
	}
	return groupSlice, nil
}

const LockGroupQuery = `SELECT id, name FROM group`

func (s *_ORM) LockGroup(ctx context.Context, queryerContext ormcommon.QueryerContext, opts ...ormopt.QueryOption) (group_.GroupSlice, error) {
	config := new(ormopt.QueryConfig)
	ormopt.WithPlaceholderGenerator(DefaultPlaceholderGenerator).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyQueryOption(config)
	}
	query, args := config.ToSQL(LockGroupQuery, 1)
	ormcommon.LoggerFromContext(ctx).Debug(query)
	rows, err := queryerContext.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("queryerContext.QueryContext: %w", s.HandleError(ctx, err))
	}
	var groupSlice group_.GroupSlice
	for rows.Next() {
		group := new(group_.Group)
		err := rows.Scan(&group.ID, &group.Name)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", s.HandleError(ctx, err))
		}
		groupSlice = append(groupSlice, group)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("rows.Close: %w", s.HandleError(ctx, err))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", s.HandleError(ctx, err))
	}
	return groupSlice, nil
}

const UpdateGroupQuery = `UPDATE group SET (name) = ($1) WHERE id = $2`

func (s *_ORM) UpdateGroup(ctx context.Context, queryerContext ormcommon.QueryerContext, group *group_.Group) error {
	ormcommon.LoggerFromContext(ctx).Debug(UpdateGroupQuery)
	_, err := queryerContext.ExecContext(ctx, UpdateGroupQuery, group.Name, group.ID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const DeleteGroupByPKQuery = `DELETE FROM group WHERE id = $1`

func (s *_ORM) DeleteGroupByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, id int) error {
	ormcommon.LoggerFromContext(ctx).Debug(DeleteGroupByPKQuery)
	_, err := queryerContext.ExecContext(ctx, DeleteGroupByPKQuery, id)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}
