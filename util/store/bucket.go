package store

import "sync"

type Bucket[K comparable, T any] struct {
	rw        sync.RWMutex
	m         map[K]T
	onDefault func(K) T
}

func NewBucket[K comparable, T any](onDefault func(K) T) *Bucket[K, T] {
	return &Bucket[K, T]{
		m:         make(map[K]T),
		onDefault: onDefault,
	}
}

func (b *Bucket[K, T]) Get(key K) T {
	b.rw.RLock()
	defer b.rw.RUnlock()
	v, ok := b.m[key]
	if !ok {
		return b.onDefault(key)
	}
	return v
}

func (b *Bucket[K, T]) Set(key K, t T) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.m[key] = t
}

func (b *Bucket[K, T]) Remove(key K) {
	b.rw.Lock()
	defer b.rw.Unlock()
	delete(b.m, key)
}
