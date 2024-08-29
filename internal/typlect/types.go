package typlect

import (
	"context"
	"reflect"
	"time"
)

type NoParam struct{}

var (
	TypeAny     = reflect.TypeOf((*any)(nil)).Elem()
	TypeErr     = reflect.TypeOf((*error)(nil)).Elem()
	TypeNoParam = reflect.TypeOf(NoParam{})
	TypeCtx     = reflect.TypeOf((*context.Context)(nil)).Elem()
	TypeTime    = reflect.TypeOf(time.Time{})
)

type (
	ValFunc  = func() any
	ValFuncE = func() (any, error)
)
