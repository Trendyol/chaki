package main

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/modules/server"
	"github.com/Trendyol/chaki/modules/swagger"
)

func main() {
	app := chaki.New()

	app.Use(
		server.Module(),
		swagger.Module(),
	)

	app.Provide(
		// Controller
		NewHelloController,

		// Service
		NewService,

		// Validator
		NewHelloTextValidator,
		NewQueryTextValidator,
	)

	err := app.Start()
	if err != nil {
		panic(err)
	}
}
