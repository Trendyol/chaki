package tx

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Trendyol/chaki/modules/orm/ormtest"
	"github.com/stretchr/testify/assert"
)

type testEntity struct{ ID int }

func (te *testEntity) TableName() string { return "test" }

func Test_txer_Transaction(t *testing.T) {
	t.Run("it should begin tx and commit when no error returned", func(t *testing.T) {
		// Given
		var (
			db, mock, err = ormtest.NewPostgresMock()
			txr           = NewTransactioner(db)
			q             = `.*`
			ctx           = context.Background()
			got           = &testEntity{}
		)

		rowval, _ := driver.Int32.ConvertValue(12)

		// config for enable mocking db transactions
		db.PrepareStmt = true

		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(rowval))
		mock.ExpectCommit()

		// When
		err = txr.Transaction(ctx, func(ctx context.Context) error {
			return From(ctx).Where(q).First(got).Error
		})

		// Then
		assert.Equal(t, 12, got.ID)
		assert.NoError(t, err)
	})

	t.Run("it should begin tx and rollback when error occurred", func(t *testing.T) {
		// Given
		var (
			db, mock, err = ormtest.NewPostgresMock()
			txr           = NewTransactioner(db)
			ctx           = context.Background()
			errTest       = errors.New("test error")
		)

		// config for enable mocking db transactions
		db.PrepareStmt = true

		assert.NoError(t, err)
		mock.ExpectBegin()
		mock.ExpectRollback()

		// When
		err = txr.Transaction(ctx, func(ctx context.Context) error {
			return errTest
		})

		// Then
		assert.Equal(t, errTest, err)
	})
}
