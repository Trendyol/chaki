package appctx

import "context"

type valuerKey string

type Valuer[T any] struct {
	key       valuerKey
	defaulter func() T
}

func NewValuer[T any](key string, defaultValue ...T) *Valuer[T] {
	vr := &Valuer[T]{key: valuerKey(key)}
	if len(defaultValue) == 1 {
		vr.defaulter = func() T { return defaultValue[0] }
	}
	return vr
}

func (v *Valuer[T]) OnDefault(fn func() T) *Valuer[T] {
	v.defaulter = fn
	return v
}

func (gt *Valuer[T]) Get(ctx context.Context, defaultVal ...T) T {
	vv := ctx.Value(gt.key)
	if vv == nil {
		if len(defaultVal) == 1 {
			return defaultVal[0]
		}
		if gt.defaulter != nil {
			return gt.defaulter()
		}
		return gt.emptyVal()
	}

	return vv.(T)
}

func (gt *Valuer[T]) Set(ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, gt.key, v)
}

func (gr *Valuer[T]) emptyVal() T {
	var t T
	return t
}
