// Code generated by ormgen; DO NOT EDIT.
//
// source: user
package user

import (
	"context"
	"strconv"

	ormcommon "github.com/hakadoriya/ormgen/examples/generated/postgres/ormcommon"
	ormopt "github.com/hakadoriya/ormgen/examples/generated/postgres/ormopt"
	user_ "github.com/hakadoriya/ormgen/examples/model/user"
)

var (
	DefaultBulkInsertMaxPlaceholdersPerQuery = 2000
	DefaultPlaceholderGenerator              = placeholderGenerator
)

var strconvItoa = strconv.Itoa

type ORM interface {
	InsertUser(ctx context.Context, queryerContext ormcommon.QueryerContext, user *user_.User) error
	BulkInsertUser(ctx context.Context, queryerContext ormcommon.QueryerContext, userSlice []*user_.User) error
	GetUserByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, user_id int) (*user_.User, error)
	GetUserByUsername(ctx context.Context, queryerContext ormcommon.QueryerContext, username string) (*user_.User, error)
	ListUserByUsernameAndAddress(ctx context.Context, queryerContext ormcommon.QueryerContext, username string, address string, opts ...ormopt.ResultOption) (user_.UserSlice, error)
	ListUser(ctx context.Context, queryerContext ormcommon.QueryerContext, opts ...ormopt.QueryOption) (user_.UserSlice, error)
	UpdateUser(ctx context.Context, queryerContext ormcommon.QueryerContext, user *user_.User) error
	DeleteUserByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, user_id int) error
	DeleteUserByUsername(ctx context.Context, queryerContext ormcommon.QueryerContext, username string) error
	DeleteUserByUsernameAndAddress(ctx context.Context, queryerContext ormcommon.QueryerContext, username string, address string) error
	InsertAdminUser(ctx context.Context, queryerContext ormcommon.QueryerContext, admin_user *user_.AdminUser) error
	BulkInsertAdminUser(ctx context.Context, queryerContext ormcommon.QueryerContext, adminUserSlice []*user_.AdminUser) error
	GetAdminUserByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, admin_user_id int) (*user_.AdminUser, error)
	ListAdminUser(ctx context.Context, queryerContext ormcommon.QueryerContext, opts ...ormopt.QueryOption) (user_.AdminUserSlice, error)
	UpdateAdminUser(ctx context.Context, queryerContext ormcommon.QueryerContext, admin_user *user_.AdminUser) error
	DeleteAdminUserByPK(ctx context.Context, queryerContext ormcommon.QueryerContext, admin_user_id int) error
}

func NewORM(opts ...ORMOption) ORM {
	o := new(_ORM)
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

type ORMOption interface {
	apply(o *_ORM)
}

func WithORMOptionHandleErrorFunc(handleErrorFunc func(ctx context.Context, err error) error) ORMOption {
	return &ormOptionHandleErrorFunc{handleErrorFunc: handleErrorFunc}
}

type ormOptionHandleErrorFunc struct {
	handleErrorFunc func(ctx context.Context, err error) error
}

func (o *ormOptionHandleErrorFunc) apply(s *_ORM) {
	s.HandleErrorFunc = o.handleErrorFunc
}

type _ORM struct {
	HandleErrorFunc func(ctx context.Context, err error) error
}

func (o *_ORM) HandleError(ctx context.Context, err error) error {
	if o.HandleErrorFunc != nil {
		return o.HandleErrorFunc(ctx, err)
	}
	return err
}

// placeholderGenerator is postgres dialect placeholder generator
func placeholderGenerator(i int) string {
	return "$" + strconvItoa(i)
}
