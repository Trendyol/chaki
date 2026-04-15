# Server Module

Use this page when the task affects HTTP controllers, routes, request parsing, validation, middlewares, or server error handling.

## Read Next

- Controller setup and route definitions: [controllers-and-routing.md](controllers-and-routing.md)
- Request tags, parsing, validation, and response status: [request-parsing-and-validation.md](request-parsing-and-validation.md)
- Middleware wiring and error handling: [middleware-and-errors.md](middleware-and-errors.md)

## Core Rules

- Register the module with `server.Module(...)`.
- Provide controllers through `app.Provide(...)`.
- Prefer Chaki's typed route handlers first and drop to raw Fiber handlers only when direct `*fiber.Ctx` control is necessary.
- Keep server concerns at the transport boundary; business logic should live in services.

## Server Config Example

```yaml
server:
  addr: ":8080"
  logging: true
  cors:
    allowedOrigins: "https://app.example.com,http://localhost:3000"
    allowCredentials: true
```

- `server.addr` defaults to `:8080`, so only set it when the app should listen somewhere else.
- Document `server.cors.allowCredentials` as a boolean even if some service configs serialize it as a string.
- Use `server.WithErrorHandler(...)` for transport error shaping and `github.com/Trendyol/chaki/modules/server/common.FiberConfigWrapper` for Fiber settings that are not exposed through YAML config.
