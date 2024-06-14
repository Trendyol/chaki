package wrapper

type (
	Wrapper[T any]            func(T) T
	WrapperWithInfo[T, R any] func(t T, info R) T
)

func Nop[T any]() Wrapper[T] {
	return func(t T) T { return t }
}

func Replacer[T any](t T) Wrapper[T] {
	return func(_ T) T { return t }
}

func Apply[T any](initial T, wrappers []Wrapper[T]) T {
	for _, wr := range wrappers {
		initial = wr(initial)
	}

	return initial
}
