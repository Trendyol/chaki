package store

import (
	csmap "github.com/mhmtszr/concurrent-swiss-map"
)

type Bucket[K comparable, T any] struct {
	m         *csmap.CsMap[K, T]
	onDefault func(K) T
}

func NewBucket[K comparable, T any](onDefault func(K) T) *Bucket[K, T] {
	m := csmap.Create[K, T]()
	return &Bucket[K, T]{
		m:         m,
		onDefault: onDefault,
	}
}

func (b *Bucket[K, T]) Get(key K) T {
	v, ok := b.m.Load(key)
	if !ok {
		return b.onDefault(key)
	}
	return v
}

func (b *Bucket[K, T]) Set(key K, t T) {
	b.m.Store(key, t)
}

func (b *Bucket[K, T]) Remove(key K) {
	b.m.Delete(key)
}

func (b *Bucket[K, T]) Has(key K) bool {
	return b.m.Has(key)
}
