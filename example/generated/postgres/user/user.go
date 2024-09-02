// Code generated by ormgen; DO NOT EDIT.
//
// source: user/user.go
package user

import (
	"context"
	"fmt"

	ormgen "github.com/hakadoriya/ormgen/example/generated/postgres/ormgen"
	user_ "github.com/hakadoriya/ormgen/example/model/user"
)

const CreateUserQuery = `INSERT INTO user (user_id, username, address, group_id) VALUES ($1, $2, $3, $4)`

func (s *_ORM) CreateUser(ctx context.Context, queryerContext ormgen.QueryerContext, user *user_.User) error {
	ormgen.LoggerFromContext(ctx).Debug(CreateUserQuery)
	_, err := queryerContext.ExecContext(ctx, CreateUserQuery, user.UserID, user.Username, user.Address, user.GroupID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const GetUserByPKQuery = `SELECT user_id, username, address, group_id FROM user WHERE user_id = $1`

func (s *_ORM) GetUserByPK(ctx context.Context, queryerContext ormgen.QueryerContext, user_id int) (*user_.User, error) {
	ormgen.LoggerFromContext(ctx).Debug(GetUserByPKQuery)
	row := queryerContext.QueryRowContext(ctx, GetUserByPKQuery, user_id)
	user := new(user_.User)
	err := row.Scan(&user.UserID, &user.Username, &user.Address, &user.GroupID)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return user, nil
}

const GetUserByUsernameQuery = `SELECT user_id, username, address, group_id FROM user WHERE username = $1`

func (s *_ORM) GetUserByUsername(ctx context.Context, queryerContext ormgen.QueryerContext, username string) (*user_.User, error) {
	ormgen.LoggerFromContext(ctx).Debug(GetUserByUsernameQuery)
	row := queryerContext.QueryRowContext(ctx, GetUserByUsernameQuery, username)
	user := new(user_.User)
	err := row.Scan(&user.UserID, &user.Username, &user.Address, &user.GroupID)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return user, nil
}

const ListUserByUsernameAndAddressQuery = `SELECT user_id, username, address, group_id FROM user WHERE (username = $1 AND address = $2)`

func (s *_ORM) ListUserByUsernameAndAddress(ctx context.Context, queryerContext ormgen.QueryerContext, username string, address string, opts ...ormgen.ResultOption) (user_.UserSlice, error) {
	config := new(ormgen.QueryConfig)
	ormgen.WithPlaceholderGenerator(ormgen.PlaceholderGeneratorMap["postgres"]).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyResultOption(config)
	}
	query, args := config.ToSQL(ListUserByUsernameAndAddressQuery, 3)
	ormgen.LoggerFromContext(ctx).Debug(query)
	rows, err := queryerContext.QueryContext(ctx, query, append([]interface{}{ username, address }, args...))
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

func (s *_ORM) ListUser(ctx context.Context, queryerContext ormgen.QueryerContext, opts ...ormgen.QueryOption) (user_.UserSlice, error) {
	config := new(ormgen.QueryConfig)
	ormgen.WithPlaceholderGenerator(ormgen.PlaceholderGeneratorMap["postgres"]).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyQueryOption(config)
	}
	query, args := config.ToSQL(ListUserQuery, 1)
	ormgen.LoggerFromContext(ctx).Debug(query)
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

const UpdateUserByPKQuery = `UPDATE user SET (username, address, group_id) = ($1, $2, $3) WHERE user_id = $4`

func (s *_ORM) UpdateUserByPK(ctx context.Context, queryerContext ormgen.QueryerContext, user *user_.User) error {
	ormgen.LoggerFromContext(ctx).Debug(UpdateUserByPKQuery)
	_, err := queryerContext.ExecContext(ctx, UpdateUserByPKQuery, user.Username, user.Address, user.GroupID, user.UserID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const CreateAdminUserQuery = `INSERT INTO admin_user (admin_user_id, username, group_id) VALUES ($1, $2, $3)`

func (s *_ORM) CreateAdminUser(ctx context.Context, queryerContext ormgen.QueryerContext, admin_user *user_.AdminUser) error {
	ormgen.LoggerFromContext(ctx).Debug(CreateAdminUserQuery)
	_, err := queryerContext.ExecContext(ctx, CreateAdminUserQuery, admin_user.AdminUserID, admin_user.Username, admin_user.GroupID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}

const GetAdminUserByPKQuery = `SELECT admin_user_id, username, group_id FROM admin_user WHERE admin_user_id = $1`

func (s *_ORM) GetAdminUserByPK(ctx context.Context, queryerContext ormgen.QueryerContext, admin_user_id int) (*user_.AdminUser, error) {
	ormgen.LoggerFromContext(ctx).Debug(GetAdminUserByPKQuery)
	row := queryerContext.QueryRowContext(ctx, GetAdminUserByPKQuery, admin_user_id)
	adminUser := new(user_.AdminUser)
	err := row.Scan(&adminUser.AdminUserID, &adminUser.Username, &adminUser.GroupID)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", s.HandleError(ctx, err))
	}
	return adminUser, nil
}

const ListAdminUserQuery = `SELECT admin_user_id, username, group_id FROM admin_user`

func (s *_ORM) ListAdminUser(ctx context.Context, queryerContext ormgen.QueryerContext, opts ...ormgen.QueryOption) (user_.AdminUserSlice, error) {
	config := new(ormgen.QueryConfig)
	ormgen.WithPlaceholderGenerator(ormgen.PlaceholderGeneratorMap["postgres"]).ApplyResultOption(config)
	for _, o := range opts {
		o.ApplyQueryOption(config)
	}
	query, args := config.ToSQL(ListAdminUserQuery, 1)
	ormgen.LoggerFromContext(ctx).Debug(query)
	rows, err := queryerContext.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("queryerContext.QueryContext: %w", s.HandleError(ctx, err))
	}
	var admin_userSlice user_.AdminUserSlice
	for rows.Next() {
		admin_user := new(user_.AdminUser)
		err := rows.Scan(&admin_user.AdminUserID, &admin_user.Username, &admin_user.GroupID)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", s.HandleError(ctx, err))
		}
		admin_userSlice = append(admin_userSlice, admin_user)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("rows.Close: %w", s.HandleError(ctx, err))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", s.HandleError(ctx, err))
	}
	return admin_userSlice, nil
}

const UpdateAdminUserByPKQuery = `UPDATE admin_user SET (username, group_id) = ($1, $2) WHERE admin_user_id = $3`

func (s *_ORM) UpdateAdminUserByPK(ctx context.Context, queryerContext ormgen.QueryerContext, admin_user *user_.AdminUser) error {
	ormgen.LoggerFromContext(ctx).Debug(UpdateAdminUserByPKQuery)
	_, err := queryerContext.ExecContext(ctx, UpdateAdminUserByPKQuery, admin_user.Username, admin_user.GroupID, admin_user.AdminUserID)
	if err != nil {
		return fmt.Errorf("queryerContext.ExecContext: %w", s.HandleError(ctx, err))
	}
	return nil
}