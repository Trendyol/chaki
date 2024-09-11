package main

import (
	"fmt"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/server"
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/route"
	"github.com/Trendyol/chaki/modules/swagger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
	)

	app.Use(
		server.Module(),
		swagger.Module(),
	)

	app.Provide(
		newController,
	)

	logger.Fatal(app.Start())
}

type greetRequest struct {
	Name string `param:"name" validate:"required"`
}

type fooController struct {
	controller.Controller
}

func newController() controller.Controller {
	return &fooController{
		Controller: controller.New("foo-controller").SetPrefix("/foo"),
	}
}

func (ct *fooController) greet(c *fiber.Ctx) error {
	return c.JSON(fmt.Sprintf("Hello %s", c.Params("name")))
}

func (ct *fooController) Routes() []route.Route {
	return []route.Route{
		route.FiberGet[greetRequest, string]("/:name", ct.greet),
	}
}
