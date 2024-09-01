package by

import "github.com/Trendyol/chaki/util/slc"

type ByFunc[T any] func(T) bool

func Value[T comparable](t T) ByFunc[T] {
	return func(tt T) bool {
		return tt == t
	}
}

func NotValue[T comparable](t T) ByFunc[T] {
	return func(tt T) bool {
		return tt != t
	}
}

func In[T ~[]E, E comparable](a T) ByFunc[E] {
	return func(e E) bool {
		return slc.Contains(a, e)
	}
}
