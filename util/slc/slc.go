package slc

import (
	"github.com/Trendyol/chaki/util/maps"
	"golang.org/x/exp/slices"
)

func Map[T, R any](a []T, f func(T) R) []R {
	res := make([]R, len(a))
	for i := range a {
		res[i] = f(a[i])
	}
	return res
}

func MapI[T, R any](a []T, f func(T, int) R) []R {
	res := make([]R, len(a))
	for i := range a {
		res[i] = f(a[i], i)
	}
	return res
}

func FlatMap[T, R any](a []T, f func(T) []R) []R {
	res := make([]R, 0)
	for i := range a {
		res = append(res, f(a[i])...)
	}
	return res
}

func ForEach[T any](a []T, f func(T)) {
	for i := range a {
		f(a[i])
	}
}

func Filter[T any](a []T, f func(T) bool) []T {
	res := make([]T, 0, len(a))
	for i := range a {
		if f(a[i]) {
			res = append(res, a[i])
		}
	}
	return res
}

func Reduce[T, R any](a []T, f func(r R, t T) R, defaultVal ...R) R {
	var r R
	if len(defaultVal) > 0 {
		r = defaultVal[0]
	}

	for i := range a {
		r = f(r, a[i])
	}

	return r
}

func ReduceErr[T any](a []T, f func(t T) error) error {
	for i := range a {
		if err := f(a[i]); err != nil {
			return err
		}
	}

	return nil
}

func Contains[T ~[]E, E comparable](a T, t E) bool {
	return slices.Contains[T, E](a, t)
}

func RemoveDuplicates[T comparable](a []T) []T {
	m := Reduce(a, func(acc map[T]bool, t T) map[T]bool {
		acc[t] = true
		return acc
	}, map[T]bool{})

	return maps.Keys(m)
}

func Reverse[T ~[]E, E any](a T) T {
	r := make(T, len(a))
	copy(r, a)
	slices.Reverse(r)
	return r
}
