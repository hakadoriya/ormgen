// Code generated by ormgen; DO NOT EDIT.
//
// source: group
package group

import (
	"context"
	"strconv"

	ormopt "github.com/hakadoriya/ormgen/example/generated/mysql/ormopt"
	group_ "github.com/hakadoriya/ormgen/example/model/group"
)

var (
	DefaultBulkInsertMaxPlaceholdersPerQuery = 2000
	DefaultPlaceholderGenerator              = placeholderGenerator
)

var strconvItoa = strconv.Itoa

type ORM interface {
	InsertGroup(ctx context.Context, queryerContext ormopt.QueryerContext, group *group_.Group) error
	BulkInsertGroup(ctx context.Context, queryerContext ormopt.QueryerContext, groupSlice []*group_.Group) error
	GetGroupByPK(ctx context.Context, queryerContext ormopt.QueryerContext, id int) (*group_.Group, error)
	ListGroup(ctx context.Context, queryerContext ormopt.QueryerContext, opts ...ormopt.QueryOption) (group_.GroupSlice, error)
	UpdateGroup(ctx context.Context, queryerContext ormopt.QueryerContext, group *group_.Group) error
	DeleteGroupByPK(ctx context.Context, queryerContext ormopt.QueryerContext, id int) error
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

// placeholderGenerator is mysql dialect placeholder generator
func placeholderGenerator(_ int) string {
	return "?"
}
