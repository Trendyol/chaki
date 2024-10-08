package main

import (
	"context"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/modules/couchbase"
	"github.com/couchbase/gocb/v2"
)

type Foo struct {
	ID  string
	Bar string
}

type FooRepository interface {
	Get(ctx context.Context, id string) (*Foo, error)
}

type fooRepository struct {
	coll *gocb.Collection
}

func NewRepository(cluster *gocb.Cluster, cfg *config.Config) FooRepository {
	return &fooRepository{
		coll: cluster.Bucket(cfg.GetString("couchbase.bucketname")).DefaultCollection(),
	}
}

func (r *fooRepository) Get(ctx context.Context, id string) (*Foo, error) {
	result := new(Foo)
	resp, err := r.coll.Get(id, &gocb.GetOptions{
		// Enable tracing
		ParentSpan: couchbase.ParentSpan(ctx),
	})
	if err != nil {
		return nil, err
	}

	if err := resp.Content(result); err != nil {
		return nil, err
	}

	return result, nil
}
