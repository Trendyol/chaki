package main

import (
	"context"

	"github.com/Trendyol/chaki/modules/newrelic"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	nrorm "github.com/Trendyol/chaki/modules/newrelic/orm"
	nrserver "github.com/Trendyol/chaki/modules/newrelic/server"

	"github.com/Trendyol/chaki/modules/orm"
	postgresdriver "github.com/Trendyol/chaki/modules/orm/driver/postgres"
	"github.com/Trendyol/chaki/modules/otel"
	otelserver "github.com/Trendyol/chaki/modules/otel/server"
	"github.com/Trendyol/chaki/modules/server"
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/route"
	"github.com/Trendyol/chaki/modules/swagger"
)

// Same example with postgres, but with newrelic.
func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
	)

	app.Use(
		otel.Module(
			otelserver.WithServer(),
		),
		orm.Module(postgresdriver.New()),
		server.Module(),
		swagger.Module(),

		newrelic.Module(
			nrserver.WithServer(),
			nrorm.WithGorm(),
		),
	)

	app.Provide(
		newFooRepository,
		newHelloController,
	)

	if err := app.Start(); err != nil {
		logger.Fatal(err)
	}
}

type getFooReq struct {
	Id int `query:"id"`
}

type HelloController struct {
	*controller.Base
	repo *fooRepository
}

func newHelloController(repo *fooRepository) controller.Controller {
	return &HelloController{
		Base: controller.New("Hello Controller").SetPrefix("/hello"),
		repo: repo,
	}
}

func (c *HelloController) Routes() []route.Route {
	return []route.Route{
		route.Get("/", c.getFoo).Name("Get Foo"),
	}
}

func (c *HelloController) getFoo(ctx context.Context, req getFooReq) (*foo, error) {
	logger.From(ctx).Info("trace id and spand id will be logged wiht message")
	return c.repo.getFoo(ctx, req.Id)
}

type foo struct {
	Id  int    `gorm:"primaryKey"`
	Bar string `gorm:"bar"`
}

type fooRepository struct {
	gp orm.GormProvider
}

func newFooRepository(gp orm.GormProvider) *fooRepository {
	return &fooRepository{gp}
}

func (repo *fooRepository) getFoo(ctx context.Context, id int) (*foo, error) {
	f := &foo{}
	err := repo.gp.Get(ctx).Where("id = ?", id).Find(f).Error
	if err != nil {
		return nil, err
	}
	return f, nil
}
