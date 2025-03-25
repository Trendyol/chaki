package orm

import (
	"context"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/orm/driver"
	"github.com/Trendyol/chaki/modules/orm/tx"
	"github.com/Trendyol/chaki/util/health"
	"gorm.io/gorm"
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
		newLivenessReadinessProbe,
	)

	m.AddProvideHook(
		module.ProvideHook{
			Wrap:  asGromWrapper.Value,
			Match: asGromWrapper.Match,
		},
	)

	return m
}

func newLivenessReadinessProbe(db *gorm.DB) health.Probe {
	f := func(ctx context.Context) error {
		var i int
		err := db.Raw("SELECT 1").Scan(&i).Error
		return err
	}

	return health.NewProbe(f, f)
}
