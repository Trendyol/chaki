package main

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/couchbase"
	"github.com/Trendyol/chaki/modules/server"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
	)

	app.Use(
		server.Module(),
		couchbase.Module(),
	)

	app.Provide(
		NewRepository,
		newController,
	)

	logger.Fatal(app.Start())
}
