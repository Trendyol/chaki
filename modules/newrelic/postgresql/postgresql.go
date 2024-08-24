package psql

import "github.com/Trendyol/chaki/config"

func WithPostgesql() any {
	return func(cfg *config.Config) *config.Config {
		if key := "gorm.use"; cfg.Exists(key) && cfg.GetString(key) == "postgres" {
			cfg.Set("gorm.postgres.dialect.drivername", "nrpgx")
		}

		return cfg
	}
}
