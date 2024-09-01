package txmock

import (
	"context"
	"database/sql"

	"github.com/Trendyol/chaki/modules/orm/tx"
)

type transactioner struct{}

func New() tx.Transactioner {
	return &transactioner{}
}

func (t *transactioner) Transaction(ctx context.Context, f func(context.Context) error, opts ...*sql.TxOptions) error {
	return f(ctx)
}
