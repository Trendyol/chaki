package di

import (
	"fmt"

	"go.uber.org/fx"
)

func ByValue[T any](t T) func() T {
	return func() T {
		return t
	}
}

func As[T any](ctr any, group string) any {
	return fx.Annotate(
		ctr,
		fx.As(new(T)),
		fx.ResultTags(fmt.Sprintf(`group:"%s"`, group)),
	)
}

func AsOnlyGroup(ctr any, group string) any {
	return fx.Annotate(
		ctr,
		fx.ResultTags(fmt.Sprintf(`group:"%s"`, group)),
	)
}

func AsHandler(ctr any, group string) any {
	return fx.Annotate(
		ctr,
		fx.ParamTags(fmt.Sprintf(`group:"%s"`, group)),
	)
}
