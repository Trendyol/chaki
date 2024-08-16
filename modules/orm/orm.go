package orm

import (
	"context"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/modules/orm/tx"
	"github.com/Trendyol/chaki/util/wrapper"
	"gorm.io/gorm"
)

type GormWrapper wrapper.Wrapper[*gorm.DB]

type GormProvider interface {
	Get(ctx context.Context) *gorm.DB
}

type provider struct {
	db *gorm.DB
}

func newGorm(cfg *config.Config, opts *options, wrappers []GormWrapper) (*gorm.DB, error) {
	db, err := opts.driver.BuildGorm(cfg)
	if err != nil {
		return nil, err
	}

	for _, wr := range wrappers {
		db = wr(db)
	}

	return db, err
}

func newGormProvider(db *gorm.DB, _ *options) GormProvider {
	return &provider{
		db: db,
	}
}

func (p *provider) Get(ctx context.Context) *gorm.DB {
	if found := tx.From(ctx); found != nil {
		return found
	}
	return p.db.WithContext(ctx)
}
