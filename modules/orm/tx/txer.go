package tx

import (
	"context"
	"database/sql"

	"github.com/Trendyol/chaki/util/appctx"
	"gorm.io/gorm"
)

var txCtxValuer = appctx.NewValuer[*gorm.DB]("tx_gorm_key")

func From(ctx context.Context) *gorm.DB {
	return txCtxValuer.Get(ctx)
}

func WithContext(parent context.Context, tx *gorm.DB) context.Context {
	return txCtxValuer.Set(parent, tx)
}

type Transactioner interface {
	Transaction(ctx context.Context, f func(ctx context.Context) error, opts ...*sql.TxOptions) error
}

type txer struct {
	db *gorm.DB
}

func NewTransactioner(db *gorm.DB) Transactioner {
	return &txer{
		db: db,
	}
}

func (txr *txer) Transaction(ctx context.Context, f func(ctx context.Context) error, opts ...*sql.TxOptions) error {
	return txr.db.Transaction(func(tx *gorm.DB) error {
		return f(txCtxValuer.Set(ctx, tx.WithContext(ctx)))
	}, opts...)
}
