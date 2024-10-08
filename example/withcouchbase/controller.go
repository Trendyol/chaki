package main

import (
	"context"

	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/route"
)

type getFooRequest struct {
	ID string `param:"id" validate:"required"`
}

type fooController struct {
	controller.Controller
	repo FooRepository
}

func newController(repo FooRepository) controller.Controller {
	return &fooController{
		Controller: controller.New("foo-controller").SetPrefix("/foo"),
		repo:       repo,
	}
}

func (ct *fooController) getFoo(ctx context.Context, req getFooRequest) (*Foo, error) {
	return ct.repo.Get(ctx, req.ID)
}

func (ct *fooController) Routes() []route.Route {
	return []route.Route{
		route.Get("/:id", ct.getFoo),
	}
}
