package typlect

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"runtime"

	"github.com/Trendyol/chaki/util/slc"
)

var (
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

func IsInterface[T any]() (ok bool) {
	defer func() {
		if msg := recover(); msg != nil {
			ok = false
		}
	}()
	t := reflect.TypeOf((*T)(nil)).Elem()
	reflect.TypeOf((*any)(nil)).Elem().Implements(t)
	ok = true
	return
}

func GetType[T any]() reflect.Type {
	var (
		v T
		t = reflect.TypeOf(v)
	)

	if t == TypeAny || t == nil {
		return TypeNoParam
	}

	return t
}

func IsFunc(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

func IsVariadicFunc(v any) bool {
	return reflect.TypeOf(v).IsVariadic()
}

func FuncInOut(v any) ([]reflect.Type, []reflect.Type) {
	var (
		in  []reflect.Type
		out []reflect.Type
		vt  = reflect.TypeOf(v)
	)

	for i := 0; i < vt.NumIn(); i++ {
		in = append(in, vt.In(i))
	}

	for i := 0; i < vt.NumOut(); i++ {
		out = append(out, vt.Out(i))
	}

	return in, out
}

func BuildCtrFromValue(t any) any {
	var (
		tv = reflect.ValueOf(t)
		tt = tv.Type()
		ft = reflect.FuncOf([]reflect.Type{}, []reflect.Type{tt}, false)
	)
	return reflect.MakeFunc(ft, func(args []reflect.Value) (results []reflect.Value) {
		return []reflect.Value{tv}
	}).Interface()
}

func InvokerFromValues(vs ...any) any {
	var rvs []reflect.Type
	for _, v := range vs {
		rvs = append(rvs, reflect.TypeOf(v))
	}
	return reflect.MakeFunc(
		reflect.FuncOf(rvs, []reflect.Type{}, false),
		func(args []reflect.Value) (results []reflect.Value) {
			return results
		}).Interface()
}

func BuildCtrInvoker(ctr any, h func(v any) error) any {
	var (
		ct  = reflect.TypeOf(ctr)
		cot = ct.Out(0)
		it  = []reflect.Type{cot}
		ot  = []reflect.Type{errType}
		ft  = reflect.FuncOf(it, ot, false)
	)

	return reflect.MakeFunc(ft, func(args []reflect.Value) (results []reflect.Value) {
		h(args[0].Interface())

		return []reflect.Value{reflect.Zero(errType)}
	}).Interface()
}

func CtrOf[T any](ctr any) bool {
	if !IsFunc(ctr) {
		return false
	}
	_, out := FuncInOut(ctr)
	tt := reflect.TypeOf((*T)(nil)).Elem()
	for _, o := range out {
		if o.Implements(tt) {
			return true
		}
	}
	return false
}

func CtrOfValue[T any](ctr any) bool {
	if !IsFunc(ctr) {
		return false
	}
	_, out := FuncInOut(ctr)
	tt := reflect.TypeOf((*T)(nil)).Elem()
	for _, o := range out {
		if o == tt {
			return true
		}
	}
	return false
}

func GetByReturns[T any](ctr any) (T, bool) {
	_, out := FuncInOut(ctr)
	tt := reflect.TypeOf((*T)(nil)).Elem()
	for _, o := range out {
		if o.Implements(tt) {
			return reflect.Zero(o).Interface().(T), true
		}
	}
	return reflect.Zero(tt).Interface().(T), false
}

func HasTypeOf[T any](p ...any) bool {
	ptypes := slc.Map(p, reflect.TypeOf)
	targett := reflect.TypeOf((*T)(nil)).Elem()
	return slc.Contains(ptypes, targett)
}

func FuncName(fn any) string {
	fnV := reflect.ValueOf(fn)
	if fnV.Kind() != reflect.Func {
		return fmt.Sprint(fn)
	}

	function := runtime.FuncForPC(fnV.Pointer()).Name()
	return fmt.Sprintf("%s()", cleanseFuncName(function))
}

func HasTag(s any, tag string) bool {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Struct {
		panic("not a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get(tag) != "" {
			return true
		}
	}

	return false
}

var vendorRegexp = regexp.MustCompile("^.*?/vendor/")

func cleanseFuncName(name string) string {
	if ue, err := url.QueryUnescape(name); err == nil {
		name = ue
	}
	return vendorRegexp.ReplaceAllString(name, "vendor/")
}
