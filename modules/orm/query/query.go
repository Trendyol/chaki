package query

import (
	"fmt"

	"github.com/Trendyol/chaki/util/wrapper"
	"gorm.io/gorm"
)

type Query wrapper.Wrapper[*gorm.DB]

func ByValue(field string, val any) Query {
	return byMatch(field, "=", val)
}

func ByIn[T any](field string, val []T) Query {
	return byMatch(field, "IN", val)
}

func ById[T any](id T) Query {
	return func(d *gorm.DB) *gorm.DB { return d.Where(id) }
}

func ByModel(m any) Query {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where(m)
	}
}

func byMatch(field string, op string, val any) Query {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where(fmt.Sprintf("%s %s ?", field, op), val)
	}
}
