// Code generated by ormgen; DO NOT EDIT.
//
// source: group/group.go
package group

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	ormcommon "github.com/hakadoriya/ormgen/examples/generated/mysql/ormcommon"
	ormopt "github.com/hakadoriya/ormgen/examples/generated/mysql/ormopt"

	group_ "github.com/hakadoriya/ormgen/examples/model/group"
)

const InsertGroupQuery = `INSERT INTO group (id, name) VALUES (?, ?)`

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
	for i := 0; i < len(groupSlice); i += maxRowsPerQuery {
		placeholderIdx := placeholderStartAt
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

const GetGroupByPKQuery = `SELECT id, name FROM group WHERE id = ?`

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

const SelectGroupForUpdateByPKQuery = `SELECT id, name FROM group WHERE id = ? FOR UPDATE`

func (s *_ORM) SelectGroupForUpdateByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, id int) (*group_.Group, error) {
	ormcommon.LoggerFromContext(ctx).Debug(SelectGroupForUpdateByPKQuery)
	row := queryerContext.QueryRowContext(ctx, SelectGroupForUpdateByPKQuery, id)
	group := new(group_.Group)
	err := row.Scan(&group.ID, &group.Name)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return group, nil
}

const GetOneGroupQuery = `SELECT id, name FROM group`

func (s *_ORM) GetOneGroup(ctx context.Context, queryerContext ormcommon.QueryerContext, opts ...ormopt.QueryOption) (*group_.Group, error) {
	config := new(ormopt.QueryConfig)
	ormopt.WithPlaceholderGenerator(DefaultPlaceholderGenerator).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyQueryOption(config)
	}
	query, args := config.ToSQL(GetOneGroupQuery, 1)

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
	if len(groupSlice) == 0 {
		return nil, fmt.Errorf("len(groupSlice)==0: %w", s.HandleError(ctx, sql.ErrNoRows))
	}
	if l := len(groupSlice); l > 1 {
		return nil, fmt.Errorf("len(groupSlice)==%d: %w", l, s.HandleError(ctx, ormcommon.ErrNotUnique))
	}
	return groupSlice[0], nil
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

const SelectGroupForUpdateQuery = `SELECT id, name FROM group`

func (s *_ORM) SelectGroupForUpdate(ctx context.Context, queryerContext ormcommon.QueryerContext, opts ...ormopt.QueryOption) (group_.GroupSlice, error) {
	config := new(ormopt.QueryConfig)
	ormopt.WithPlaceholderGenerator(DefaultPlaceholderGenerator).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyQueryOption(config)
	}
	query, args := config.ToSQL(SelectGroupForUpdateQuery, 1)
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

const UpdateGroupQuery = `UPDATE group SET (name) = (?) WHERE id = ?`

func (s *_ORM) UpdateGroup(ctx context.Context, queryerContext ormcommon.QueryerContext, group *group_.Group) error {
	ormcommon.LoggerFromContext(ctx).Debug(UpdateGroupQuery)
	_, err := queryerContext.ExecContext(ctx, UpdateGroupQuery, group.Name, group.ID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const DeleteGroupByPKQuery = `DELETE FROM group WHERE id = ?`

func (s *_ORM) DeleteGroupByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, id int) error {
	ormcommon.LoggerFromContext(ctx).Debug(DeleteGroupByPKQuery)
	_, err := queryerContext.ExecContext(ctx, DeleteGroupByPKQuery, id)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}
