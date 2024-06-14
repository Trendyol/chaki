package as

import (
	"github.com/Trendyol/chaki/internal/di"
	"github.com/Trendyol/chaki/internal/typlect"
)

// Asser is a representation of type checker for constructors.
// check value the  constructor provides and can group constructors which provide the same type
type Asser interface {
	Value(any) any
	Handler(any) any
	Match(any) bool
	Grouper() any
}

type asser[T any] struct {
	key string
	g   bool
}

// Interface is a constructor method for an Asser that can be used for Interterface types.
// not for primitive types, function types, or struct types.
func Interface[T any](key string) Asser {
	return &asser[T]{key, false}
}

// Struct is a constructor method for an Asser that can be used for primitive types, function types, or struct types.
func Struct[T any](key string) Asser {
	return &asser[T]{key, true}
}

// Value wraps the constructor with injector annotation so that it can be handled as a group
func (a *asser[T]) Value(ctr any) any {
	if a.g {
		return di.AsOnlyGroup(ctr, a.key)
	}
	return di.As[T](ctr, a.key)
}

// Handler wraps constructor with handler injector anotation.
func (a *asser[T]) Handler(ctr any) any {
	return di.AsHandler(ctr, a.key)
}

// Match checks the constructor provided value match with the Asser value
func (a *asser[T]) Match(ctr any) bool {
	if a.g {
		return typlect.CtrOfValue[T](ctr)
	}
	return typlect.CtrOf[T](ctr)
}

// Grouper returns a specific handler annotation, which groups the values as a list and then returns it.
// After providing the grouper to the provided pool of the di tool.
// So that list can be used by any constructor without any need to group annotation
func (a *asser[T]) Grouper() any {
	return a.Handler(func(args []T) []T {
		return args
	})
}
