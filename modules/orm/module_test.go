package orm

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Trendyol/chaki/modules/orm/ormtest"
	"github.com/stretchr/testify/assert"
)

func Test_newLivenessReadinessProbe(t *testing.T) {
	// Given
	db, mock, err := ormtest.NewPostgresMock()
	assert.NoError(t, err)
	probe := newLivenessReadinessProbe(db)
	query := "SELECT 1"
	rowValue, _ := driver.Int32.ConvertValue(1)
	rows := sqlmock.NewRows([]string{"1"}).AddRow(rowValue)

	mock.ExpectQuery(query).WillReturnRows(rows) // for readiness
	mock.ExpectQuery(query).WillReturnRows(rows) // for liveness

	// When
	errReadiness := probe.Readiness(context.Background())
	errLiveness := probe.Liveness(context.Background())

	// Then
	assert.NoError(t, errReadiness)
	assert.NoError(t, errLiveness)
}
