// Code generated by ormgen; DO NOT EDIT.
//
// source: gourp
package group

import (
	"context"

	ormgen "github.com/hakadoriya/ormgen/internal/lang/go/generator/test/generated/ormgen"
	gourp_ "github.com/hakadoriya/ormgen/internal/lang/go/source/test/gourp"
)

type ORM interface {
	CreateGroup(ctx context.Context, queryerContext ormgen.QueryerContext, group *gourp_.Group) error
	GetGroupByPK(ctx context.Context, queryerContext ormgen.QueryerContext, id int) (*gourp_.Group, error)
	ListGroup(ctx context.Context, queryerContext ormgen.QueryerContext, opts ...ormgen.QueryOption) (gourp_.GroupSlice, error)
	UpdateGroupByPK(ctx context.Context, queryerContext ormgen.QueryerContext, group *gourp_.Group) error
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