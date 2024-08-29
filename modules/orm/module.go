package orm

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/orm/driver"
	"github.com/Trendyol/chaki/modules/orm/tx"
)

const ModuleName = "chaki-orm-module"

var asGromWrapper = as.Struct[GormWrapper]("gormwrappers")

func Module(d driver.Driver, opts ...Option) *module.Module {
	m := module.New(ModuleName)
	o := buildOptions(opts...)
	o.driver = d

	m.Provide(
		chaki.Valuer(o),
		asGromWrapper.Grouper(),
		newGorm,
		newGormProvider,
		tx.NewTransactioner,
	)

	m.AddProvideHook(
		module.ProvideHook{
			Wrap:  asGromWrapper.Value,
			Match: asGromWrapper.Match,
		},
	)

	return m
}
