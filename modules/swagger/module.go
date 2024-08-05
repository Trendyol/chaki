package swagger

import (
	"github.com/Trendyol/chaki/module"
)

func Module() *module.Module {
	m := module.New("swagger")

	m.Provide(
		buildDocs,
		fiberWrapper,
	)

	return m
}
