package driver

import (
	"github.com/Trendyol/chaki/config"
	"gorm.io/gorm"
)

type Driver interface {
	BuildGorm(cfg *config.Config) (*gorm.DB, error)
}
