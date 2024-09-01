package repository

import (
	"context"
	"strings"

	"github.com/Trendyol/chaki/modules/orm"
	"github.com/Trendyol/chaki/modules/orm/query"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IDType interface {
	comparable
}

var ErrRecordNotFound = gorm.ErrRecordNotFound

type Repository[Id IDType, T any] interface {
	FindByID(ctx context.Context, id Id) (*T, error)
	FindOne(ctx context.Context, q query.Query) (*T, error)
	FindAll(ctx context.Context, q query.Query) ([]*T, error)
	Update(ctx context.Context, q query.Query, update any) error
	Save(ctx context.Context, t *T) (*T, error)
	SaveAll(ctx context.Context, ts []*T) error
	DeleteByID(ctx context.Context, id Id) error
	Delete(ctx context.Context, q query.Query) error
	ListPageable(ctx context.Context, q query.Query, req PageableRequest) (*PageableResponse[*T], error)
	ParseQuery(ctx context.Context, q query.Query) *gorm.DB
	Context(ctx context.Context) *gorm.DB
}

type repository[Id IDType, T any] struct {
	gp orm.GormProvider
}

func New[Id IDType, T any](gp orm.GormProvider) Repository[Id, T] {
	return &repository[Id, T]{
		gp: gp,
	}
}

func (r *repository[Id, T]) FindByID(ctx context.Context, id Id) (*T, error) {
	return r.FindOne(ctx, query.ByID(id))
}

func (r *repository[Id, T]) FindOne(ctx context.Context, q query.Query) (*T, error) {
	t := new(T)
	err := r.ParseQuery(ctx, q).First(t).Error
	return t, err
}

func (r *repository[Id, T]) FindAll(ctx context.Context, q query.Query) ([]*T, error) {
	a := make([]*T, 0)
	err := r.ParseQuery(ctx, q).Find(&a).Error
	return a, err
}

func (r *repository[Id, T]) Update(ctx context.Context, q query.Query, update any) error {
	return r.ParseQuery(ctx, q).Updates(update).Error
}

func (r *repository[Id, T]) Save(ctx context.Context, t *T) (*T, error) {
	err := r.Context(ctx).Clauses(clause.Returning{}).Save(t).Error
	return t, err
}

func (r *repository[Id, T]) SaveAll(ctx context.Context, ts []*T) error {
	return r.Context(ctx).Save(ts).Error
}

func (r *repository[Id, T]) DeleteByID(ctx context.Context, id Id) error {
	return r.Delete(ctx, query.ByID(id))
}

func (r *repository[Id, T]) Delete(ctx context.Context, q query.Query) error {
	return r.ParseQuery(ctx, q).Delete(new(T)).Error
}

func (r *repository[Id, T]) ListPageable(ctx context.Context, q query.Query, req PageableRequest) (*PageableResponse[*T], error) {
	var (
		offset = req.Page * req.Size
		sort   = strings.Join(req.Sort, ", ")
		countq = r.ParseQuery(ctx, q)
		resq   = r.ParseQuery(ctx, q).Offset(offset).Limit(req.Size)
		found  = make([]*T, 0, req.Size)
		count  int64
		eg     errgroup.Group
	)

	if sort != "" {
		resq.Order(sort)
	}

	eg.Go(func() error {
		return resq.Find(&found).Error
	})

	eg.Go(func() error {
		return countq.Count(&count).Error
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &PageableResponse[*T]{
		Page:       req.Page,
		Size:       req.Size,
		Data:       found,
		TotalCount: int(count),
	}, nil
}

func (r *repository[Id, T]) ParseQuery(ctx context.Context, q query.Query) *gorm.DB {
	return q(r.Context(ctx))
}

func (r *repository[Id, T]) Context(ctx context.Context) *gorm.DB {
	return r.gp.Get(ctx).Model(new(T))
}
