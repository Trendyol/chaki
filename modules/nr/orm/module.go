package nrorm

import (
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/orm"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const ModuleName = "chaki-nr-orm-module"

var asSegmentHook = as.Struct[SegmentHook]("nrormsegmenthooks")

func Module() *module.Module {
	return module.New(ModuleName).Provide(newGormWrapper).AddAsser(asSegmentHook)
}

func newGormWrapper(nr *newrelic.Application, hooks []SegmentHook) orm.GormWrapper {
	return func(d *gorm.DB) *gorm.DB {
		if err := d.Use(newGormPlugin(nr, hooks)); err != nil {
			logger.New().Fatal("gorm newrelic plugin can not implemented", zap.Error(err))
		}
		return d
	}
}
