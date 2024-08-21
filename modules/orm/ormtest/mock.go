package ormtest

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMock(dialectorGenerator func(db *sql.DB) gorm.Dialector) (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(dialectorGenerator(db))
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func NewPostgresMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	return NewMock(func(db *sql.DB) gorm.Dialector {
		return postgres.New(postgres.Config{
			Conn: db,
		})
	})
}
