package repository

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Trendyol/chaki/modules/orm"
	"github.com/Trendyol/chaki/modules/orm/ormtest"
	"github.com/Trendyol/chaki/modules/orm/query"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testGormProvider struct{ db *gorm.DB }

func (p *testGormProvider) Get(context.Context) *gorm.DB { return p.db }

type testModel struct {
	Id  int    `gorm:"primaryKey"`
	Foo string `gorm:"foo"`
}

// gorm method for getting table name
func (tm *testModel) TableName() string {
	return "test_table"
}

func newTestRepository() (Repository[int, testModel], sqlmock.Sqlmock) {
	var (
		db, mock, _                  = ormtest.NewPostgresMock()
		gp          orm.GormProvider = &testGormProvider{db}
		repo                         = New[int, testModel](gp)
	)
	return repo, mock
}

func buildRows(rows ...*testModel) *sqlmock.Rows {
	r := sqlmock.NewRows([]string{"id", "foo"})
	for _, row := range rows {
		idv, _ := driver.Int32.ConvertValue(row.Id)
		foov, _ := driver.String.ConvertValue(row.Foo)
		r.AddRow(idv, foov)
	}
	return r
}

func Test_respository_FindById(t *testing.T) {
	t.Run("it should run find by id query when valid requested", func(t *testing.T) {
		// Given
		repo, mock := newTestRepository()
		m := &testModel{
			Id:  12,
			Foo: "test",
		}
		mock.
			ExpectQuery(`SELECT \* FROM "test_table" WHERE "test_table"\."id" = \$1`).
			WithArgs(12, 1). // id and limit statement arguments
			WillReturnRows(buildRows(m))

		// When
		got, err := repo.FindById(context.Background(), 12)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, m, got)

	})
}

func Test_respository_FindOne(t *testing.T) {
	t.Run("it should run find one query when valid requested", func(t *testing.T) {
		// Given
		repo, mock := newTestRepository()
		m := &testModel{
			Id:  12,
			Foo: "test",
		}

		testQ := query.Query(func(d *gorm.DB) *gorm.DB {
			return d.Where(&testModel{
				Foo: "test",
			})
		})

		mock.
			ExpectQuery(`SELECT \* FROM "test_table" WHERE "test_table"\."foo" = \$1`).
			WithArgs("test", 1). // id and limit statement arguments
			WillReturnRows(buildRows(m))

		// When
		got, err := repo.FindOne(context.Background(), testQ)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, m, got)

	})
}

func Test_respository_FindAll(t *testing.T) {
	t.Run("it should run find all  when valid requested", func(t *testing.T) {
		// Given
		repo, mock := newTestRepository()
		testModels := []*testModel{
			{12, "test"},
			{35, "test"},
		}

		testQ := query.Query(func(d *gorm.DB) *gorm.DB {
			return d.Where(&testModel{
				Foo: "test",
			})
		})

		mock.
			ExpectQuery(`SELECT \* FROM "test_table" WHERE "test_table"\."foo" = \$1`).
			WithArgs("test").
			WillReturnRows(buildRows(testModels...))

		// When
		got, err := repo.FindAll(context.Background(), testQ)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, testModels, got)

	})
}

func Test_respository_Update(t *testing.T) {
	t.Run("it should run update query  when valid requested", func(t *testing.T) {
		// Given
		repo, mock := newTestRepository()

		u := &testModel{
			Foo: "test",
		}

		mock.ExpectBegin()
		mock.
			ExpectExec(`UPDATE "test_table" SET "foo"=\$1 WHERE "test_table"\."id" = \$2`).
			WithArgs("test", 12).
			WillReturnResult(driver.RowsAffected(1))
		mock.ExpectCommit()

		// When
		err := repo.Update(context.Background(), query.ById(12), u)

		// Then
		assert.NoError(t, err)

	})
}

func Test_respository_SaveAll(t *testing.T) {
	t.Run("it should run save all query when valid requested", func(t *testing.T) {
		// Given
		repo, mock := newTestRepository()

		u := []*testModel{
			{
				Id:  12,
				Foo: "test",
			},
			{
				Id:  14,
				Foo: "test_2",
			},
		}

		mock.ExpectBegin()
		mock.
			ExpectQuery(`INSERT INTO "test_table" \("foo","id"\) VALUES \(\$1,\$2\),\(\$3,\$4\) ON CONFLICT \("id"\) DO UPDATE SET "foo"="excluded"\."foo"`).
			WithArgs("test", 12, "test_2", 14).
			WillReturnRows(buildRows(u...))
		mock.ExpectCommit()

		// When
		err := repo.SaveAll(context.Background(), u)

		// Then
		assert.NoError(t, err)

	})
}

func Test_respository_Delete(t *testing.T) {
	t.Run("it should run delete query  when valid requested", func(t *testing.T) {
		// Given
		repo, mock := newTestRepository()

		q := query.ById(12)

		mock.ExpectBegin()
		mock.
			ExpectExec(`DELETE FROM "test_table" WHERE "test_table"\."id" = \$1`).
			WithArgs(12).
			WillReturnResult(driver.RowsAffected(1))

		mock.ExpectCommit()

		// When
		err := repo.Delete(context.Background(), q)

		// Then
		assert.NoError(t, err)

	})
}
