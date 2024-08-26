package psql

import (
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	nrmodule "github.com/Trendyol/chaki/modules/newrelic"
)

func WithPostgesql() nrmodule.Option {
	return nrmodule.WithSubModule(
		module.NewSubModule().Provide(psqlConfigWrapper),
	)
}

func psqlConfigWrapper() config.Wrapper {
	return func(cfg *config.Config) *config.Config {
		if key := "gorm.use"; cfg.Exists(key) && cfg.GetString(key) == "postgres" {
			cfg.Set("gorm.postgres.dialect.drivername", "nrpgx")
		}
		return cfg
	}
}
