package postgresdriver

import (
	"fmt"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/modules/orm/driver"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDriver struct{}

func New() driver.Driver {
	return &postgresDriver{}
}

func (d *postgresDriver) BuildGorm(cfg *config.Config) (*gorm.DB, error) {
	ormCfg := cfg.Of("gorm")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		ormCfg.GetString("host"),
		ormCfg.GetString("user"),
		ormCfg.GetString("password"),
		ormCfg.GetString("dbname"),
		ormCfg.GetInt("port"),
	)

	pc := postgres.Config{
		DSN: dsn,
	}

	if key := "dialect.drivername"; ormCfg.Exists(key) {
		pc.DriverName = ormCfg.GetString(key)
	}

	return gorm.Open(postgres.New(pc))
}
