package nrorm

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Trendyol/chaki/modules/orm/ormtest"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
)

type recorderHook struct {
	seg *newrelic.DatastoreSegment
}

func (h *recorderHook) hook(seg *newrelic.DatastoreSegment) *newrelic.DatastoreSegment {
	h.seg = seg
	return seg
}

func Test_gormPlugin_Initialize(t *testing.T) {
	t.Run("it should create segment when query executed", func(t *testing.T) {
		// Given
		var (
			db, mock, _ = ormtest.NewPostgresMock()
			nr, _       = newrelic.NewApplication()
			hook        = &recorderHook{nil}
			plugin      = newGormPlugin(nr, []SegmentHook{hook.hook})
			query       = "SELECT * FROM test"
		)
		pErr := db.Use(plugin)

		v, _ := driver.Bool.ConvertValue(true)
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"test"}).AddRow(v))

		// When
		gotData := make([]map[string]any, 0)
		err := db.Raw(query).First(&gotData).Error
		got := hook.seg

		// Then
		assert.NotNil(t, got)
		assert.NoError(t, pErr)
		assert.NoError(t, err)
		assert.Equal(t, got.DatabaseName, "postgres")
	})
}
