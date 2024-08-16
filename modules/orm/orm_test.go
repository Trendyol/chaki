package orm

import (
	"context"
	"testing"

	"github.com/Trendyol/chaki/modules/orm/ormtest"
	"github.com/Trendyol/chaki/modules/orm/tx"
	"github.com/stretchr/testify/assert"
)

func Test_provider_Get(t *testing.T) {
	t.Run("it should get default gorm instance when context has no tx", func(t *testing.T) {
		// Given
		var (
			db, _, err = ormtest.NewPostgresMock()
			ctx        = context.Background()
			gp         = newGormProvider(db, &options{})
		)

		// When
		got := gp.Get(ctx)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("it should get gorm instance from when context has tx", func(t *testing.T) {
		// Given
		var (
			db, _, err = ormtest.NewPostgresMock()
			txdb       = db.Begin()
			ctx        = tx.WithContext(context.Background(), txdb)
			gp         = newGormProvider(db, &options{})
		)

		defer txdb.Rollback()

		// When
		got := gp.Get(ctx)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, txdb, got)
	})
}
