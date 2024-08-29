// Code generated by ormgen; DO NOT EDIT.
//
// source: user/user.go
package user

import (
	"context"
	"fmt"
	"strings"

	ormopt "github.com/hakadoriya/ormgen/examples/generated/spanner/ormopt"

	user_ "github.com/hakadoriya/ormgen/examples/model/user"
)

const InsertUserQuery = `INSERT INTO user (user_id, username, address, group_id) VALUES (?, ?, ?, ?)`

func (s *_ORM) InsertUser(ctx context.Context, queryerContext ormopt.QueryerContext, user *user_.User) error {
	ormopt.LoggerFromContext(ctx).Debug(InsertUserQuery)
	_, err := queryerContext.ExecContext(ctx, InsertUserQuery, user.UserID, user.Username, user.Address, user.GroupID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const BulkInsertUserQueryPrefix = `INSERT INTO user (user_id, username, address, group_id) VALUES `

func (s *_ORM) BulkInsertUser(ctx context.Context, queryerContext ormopt.QueryerContext, userSlice []*user_.User) error {
	if len(userSlice) == 0 {
		return nil
	}

	// Calculate the number of placeholders per row and the maximum number of rows per query
	const (
		placeholderStartAt = 1
		placeholdersPerRow = 4
	)

	maxRowsPerQuery := DefaultBulkInsertMaxPlaceholdersPerQuery / placeholdersPerRow
	placeholderIdx := placeholderStartAt
	for i := 0; i < len(userSlice); i += maxRowsPerQuery {
		end := i + maxRowsPerQuery
		if end > len(userSlice) {
			end = len(userSlice)
		}

		chunk := userSlice[i:end]
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
			args = append(args, chunk[j].UserID, chunk[j].Username, chunk[j].Address, chunk[j].GroupID)
		}

		query := BulkInsertUserQueryPrefix + strings.Join(placeholders, ", ")
		ormopt.LoggerFromContext(ctx).Debug(query)

		_, err := queryerContext.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
		}
	}

	return nil
}

const GetUserByPKQuery = `SELECT user_id, username, address, group_id FROM user WHERE user_id = ?`

func (s *_ORM) GetUserByPK(ctx context.Context, queryerContext ormopt.QueryerContext, user_id int) (*user_.User, error) {
	ormopt.LoggerFromContext(ctx).Debug(GetUserByPKQuery)
	row := queryerContext.QueryRowContext(ctx, GetUserByPKQuery, user_id)
	user := new(user_.User)
	err := row.Scan(&user.UserID, &user.Username, &user.Address, &user.GroupID)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return user, nil
}

const GetUserByUsernameQuery = `SELECT user_id, username, address, group_id FROM user WHERE username = ?`

func (s *_ORM) GetUserByUsername(ctx context.Context, queryerContext ormopt.QueryerContext, username string) (*user_.User, error) {
	ormopt.LoggerFromContext(ctx).Debug(GetUserByUsernameQuery)
	row := queryerContext.QueryRowContext(ctx, GetUserByUsernameQuery, username)
	user := new(user_.User)
	err := row.Scan(&user.UserID, &user.Username, &user.Address, &user.GroupID)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return user, nil
}

const ListUserByUsernameAndAddressQuery = `SELECT user_id, username, address, group_id FROM user WHERE (username = ? AND address = ?)`

func (s *_ORM) ListUserByUsernameAndAddress(ctx context.Context, queryerContext ormopt.QueryerContext, username string, address string, opts ...ormopt.ResultOption) (user_.UserSlice, error) {
	config := new(ormopt.QueryConfig)
	ormopt.WithPlaceholderGenerator(DefaultPlaceholderGenerator).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyResultOption(config)
	}
	query, args := config.ToSQL(ListUserByUsernameAndAddressQuery, 3)
	ormopt.LoggerFromContext(ctx).Debug(query)
	rows, err := queryerContext.QueryContext(ctx, query, append([]interface{}{username, address}, args...)...)
	if err != nil {
		return nil, fmt.Errorf("queryerContext.QueryContext: %w", s.HandleError(ctx, err))
	}
	var userSlice user_.UserSlice
	for rows.Next() {
		user := new(user_.User)
		err := rows.Scan(&user.UserID, &user.Username, &user.Address, &user.GroupID)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", s.HandleError(ctx, err))
		}
		userSlice = append(userSlice, user)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("rows.Close: %w", s.HandleError(ctx, err))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", s.HandleError(ctx, err))
	}
	return userSlice, nil
}

const ListUserQuery = `SELECT user_id, username, address, group_id FROM user`

func (s *_ORM) ListUser(ctx context.Context, queryerContext ormopt.QueryerContext, opts ...ormopt.QueryOption) (user_.UserSlice, error) {
	config := new(ormopt.QueryConfig)
	ormopt.WithPlaceholderGenerator(DefaultPlaceholderGenerator).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyQueryOption(config)
	}
	query, args := config.ToSQL(ListUserQuery, 1)
	ormopt.LoggerFromContext(ctx).Debug(query)
	rows, err := queryerContext.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("queryerContext.QueryContext: %w", s.HandleError(ctx, err))
	}
	var userSlice user_.UserSlice
	for rows.Next() {
		user := new(user_.User)
		err := rows.Scan(&user.UserID, &user.Username, &user.Address, &user.GroupID)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", s.HandleError(ctx, err))
		}
		userSlice = append(userSlice, user)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("rows.Close: %w", s.HandleError(ctx, err))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", s.HandleError(ctx, err))
	}
	return userSlice, nil
}

const UpdateUserQuery = `UPDATE user SET (username, address, group_id) = (?, ?, ?) WHERE user_id = ?`

func (s *_ORM) UpdateUser(ctx context.Context, queryerContext ormopt.QueryerContext, user *user_.User) error {
	ormopt.LoggerFromContext(ctx).Debug(UpdateUserQuery)
	_, err := queryerContext.ExecContext(ctx, UpdateUserQuery, user.Username, user.Address, user.GroupID, user.UserID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const DeleteUserByPKQuery = `DELETE FROM user WHERE user_id = ?`

func (s *_ORM) DeleteUserByPK(ctx context.Context, queryerContext ormopt.QueryerContext, user_id int) error {
	ormopt.LoggerFromContext(ctx).Debug(DeleteUserByPKQuery)
	_, err := queryerContext.ExecContext(ctx, DeleteUserByPKQuery, user_id)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const DeleteUserByUsernameQuery = `DELETE FROM user WHERE username = ?`

func (s *_ORM) DeleteUserByUsername(ctx context.Context, queryerContext ormopt.QueryerContext, username string) error {
	ormopt.LoggerFromContext(ctx).Debug(DeleteUserByUsernameQuery)
	_, err := queryerContext.ExecContext(ctx, DeleteUserByUsernameQuery, username)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const DeleteUserByUsernameAndAddressQuery = `DELETE FROM user WHERE (username = ? AND address = ?)`

func (s *_ORM) DeleteUserByUsernameAndAddress(ctx context.Context, queryerContext ormopt.QueryerContext, username string, address string) error {
	ormopt.LoggerFromContext(ctx).Debug(DeleteUserByUsernameAndAddressQuery)
	_, err := queryerContext.ExecContext(ctx, DeleteUserByUsernameAndAddressQuery, username, address)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const InsertAdminUserQuery = `INSERT INTO admin_user (admin_user_id, username, group_id) VALUES (?, ?, ?)`

func (s *_ORM) InsertAdminUser(ctx context.Context, queryerContext ormopt.QueryerContext, admin_user *user_.AdminUser) error {
	ormopt.LoggerFromContext(ctx).Debug(InsertAdminUserQuery)
	_, err := queryerContext.ExecContext(ctx, InsertAdminUserQuery, admin_user.AdminUserID, admin_user.Username, admin_user.GroupID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const BulkInsertAdminUserQueryPrefix = `INSERT INTO admin_user (admin_user_id, username, group_id) VALUES `

func (s *_ORM) BulkInsertAdminUser(ctx context.Context, queryerContext ormopt.QueryerContext, adminUserSlice []*user_.AdminUser) error {
	if len(adminUserSlice) == 0 {
		return nil
	}

	// Calculate the number of placeholders per row and the maximum number of rows per query
	const (
		placeholderStartAt = 1
		placeholdersPerRow = 3
	)

	maxRowsPerQuery := DefaultBulkInsertMaxPlaceholdersPerQuery / placeholdersPerRow
	placeholderIdx := placeholderStartAt
	for i := 0; i < len(adminUserSlice); i += maxRowsPerQuery {
		end := i + maxRowsPerQuery
		if end > len(adminUserSlice) {
			end = len(adminUserSlice)
		}

		chunk := adminUserSlice[i:end]
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
			args = append(args, chunk[j].AdminUserID, chunk[j].Username, chunk[j].GroupID)
		}

		query := BulkInsertAdminUserQueryPrefix + strings.Join(placeholders, ", ")
		ormopt.LoggerFromContext(ctx).Debug(query)

		_, err := queryerContext.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
		}
	}

	return nil
}

const GetAdminUserByPKQuery = `SELECT admin_user_id, username, group_id FROM admin_user WHERE admin_user_id = ?`

func (s *_ORM) GetAdminUserByPK(ctx context.Context, queryerContext ormopt.QueryerContext, admin_user_id int) (*user_.AdminUser, error) {
	ormopt.LoggerFromContext(ctx).Debug(GetAdminUserByPKQuery)
	row := queryerContext.QueryRowContext(ctx, GetAdminUserByPKQuery, admin_user_id)
	adminUser := new(user_.AdminUser)
	err := row.Scan(&adminUser.AdminUserID, &adminUser.Username, &adminUser.GroupID)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return adminUser, nil
}

const ListAdminUserQuery = `SELECT admin_user_id, username, group_id FROM admin_user`

func (s *_ORM) ListAdminUser(ctx context.Context, queryerContext ormopt.QueryerContext, opts ...ormopt.QueryOption) (user_.AdminUserSlice, error) {
	config := new(ormopt.QueryConfig)
	ormopt.WithPlaceholderGenerator(DefaultPlaceholderGenerator).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyQueryOption(config)
	}
	query, args := config.ToSQL(ListAdminUserQuery, 1)
	ormopt.LoggerFromContext(ctx).Debug(query)
	rows, err := queryerContext.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("queryerContext.QueryContext: %w", s.HandleError(ctx, err))
	}
	var adminUserSlice user_.AdminUserSlice
	for rows.Next() {
		admin_user := new(user_.AdminUser)
		err := rows.Scan(&admin_user.AdminUserID, &admin_user.Username, &admin_user.GroupID)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", s.HandleError(ctx, err))
		}
		adminUserSlice = append(adminUserSlice, admin_user)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("rows.Close: %w", s.HandleError(ctx, err))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", s.HandleError(ctx, err))
	}
	return adminUserSlice, nil
}

const UpdateAdminUserQuery = `UPDATE admin_user SET (username, group_id) = (?, ?) WHERE admin_user_id = ?`

func (s *_ORM) UpdateAdminUser(ctx context.Context, queryerContext ormopt.QueryerContext, admin_user *user_.AdminUser) error {
	ormopt.LoggerFromContext(ctx).Debug(UpdateAdminUserQuery)
	_, err := queryerContext.ExecContext(ctx, UpdateAdminUserQuery, admin_user.Username, admin_user.GroupID, admin_user.AdminUserID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const DeleteAdminUserByPKQuery = `DELETE FROM admin_user WHERE admin_user_id = ?`

func (s *_ORM) DeleteAdminUserByPK(ctx context.Context, queryerContext ormopt.QueryerContext, admin_user_id int) error {
	ormopt.LoggerFromContext(ctx).Debug(DeleteAdminUserByPKQuery)
	_, err := queryerContext.ExecContext(ctx, DeleteAdminUserByPKQuery, admin_user_id)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}
