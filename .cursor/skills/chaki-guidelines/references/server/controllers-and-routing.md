# Controllers And Routing

Use this page when defining or changing HTTP endpoints.

## Preferred Pattern

- Implement `controller.Controller`.
- Embed `*controller.Base` in the controller type.
- Build the base with `controller.New(controllerName).SetPrefix(prefix)`.
- Return routes from `Routes() []route.Route`.

## Route Rules

- Prefer typed handlers with the shape `func(context.Context, Req) (Res, error)`.
- Use `route.Get`, `route.Post`, `route.Put`, and similar helpers for standard route definitions.
- Name routes so they remain identifiable in middleware, Swagger, and diagnostics.
- Use `route.NoParam` when a route does not require an input payload.

## When To Use Raw Fiber Routes

- Use `route.FiberGet`, `route.FiberPost`, and similar helpers only when direct `*fiber.Ctx` access is required.
- If raw Fiber handlers are used, keep parsing and response shaping localized and avoid moving business rules into the controller.

## Design Guidance

- Keep controller names, prefixes, and other reused route identifiers in constants when they appear in multiple places.
- Keep controllers thin: translate transport input to service calls and return domain results or mapped errors.
- If Fiber itself needs extra configuration, provide a `common.FiberConfigWrapper` rather than forking the module pattern.
