package chaki

import (
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/module"
	"go.uber.org/fx"
)

func configModule(opts configOptions) (*module.Module, error) {
	m := module.New("config")
	if opts.disabled {
		return m, nil
	}

	cfg, err := config.NewConfigFromPaths(opts.path, opts.referencePaths)
	if err != nil {
		return nil, err
	}

	m.Provide(Valuer(cfg))

	return m, nil
}

func loggerModule() *module.Module {
	syncLogger := func(lc fx.Lifecycle) {
		lc.Append(fx.StartStopHook(logger.Init, logger.Sync))
	}
	return module.New("logger").Invoke(syncLogger)
}
