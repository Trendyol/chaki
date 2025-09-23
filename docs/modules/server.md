# Server Module Documentation

## Introduction

The Server module in Chaki provides a Fiber-based HTTP server framework. It allows you to build web servers with controllers, routes, middlewares, validation, and response handling. It integrates seamlessly with other Chaki modules like Swagger, OTEL, New Relic, and ORM.

Key features:
- Controller-based routing
- Automatic request parsing (query, params, body, headers, cookies)
- Built-in validation
- Middleware support
- Health checks and logging

## Setup

To use the Server module:

1. Import the module:
   ```go
   import "github.com/Trendyol/chaki/modules/server"
   ```

2. Add it to your Chaki application:
   ```go
   app := chaki.New()
   app.Use(server.Module(/* options */))
   ```

3. Provide controllers:
   ```go
   app.Provide(NewYourController)
   ```

4. Start the app:
   ```go
   app.Start()
   ```

Optional: Pass options like `server.WithErrorHandler(customHandler)`.

## Configuration

Configure via YAML or code. Default config under `server` key:

- `addr`: Server address (default: ":8080")
- `bodylimit`: Max body size (default: 4MB)
- `readbuffersize`: Read buffer size (default: 16KB)
- `readtimeout`: Read timeout (default: 10s)
- `writetimeout`: Write timeout (default: 10s)
- `strictrouting`: Enable strict routing (default: false)
- `logging`: Enable request logging (default: false)
- `cors`: CORS settings (e.g., `allowedOrigins`, `allowCredentials`)
- `healthcheck.endpoints.liveness`: Liveness probe path (default: "/__monitor/live")
- `healthcheck.endpoints.readiness`: Readiness probe path (default: "/__monitor/ready")

Example config.yaml:
```yaml
server:
  addr: ":8080"
  logging: true
  cors:
    allowedOrigins: "http://example.com,http://anotherexample.com"
    allowCredentials: "true"
```

This configuration allows CORS requests from specific origins and enables credentials.

### Advanced Fiber Configuration

If you need to override Fiber settings not exposed in the standard configuration, provide a `common.FiberConfigWrapper`:

```go
import "github.com/Trendyol/chaki/modules/server/common"

func fiberWrapper() common.FiberConfigWrapper {
    return func(cfg fiber.Config) fiber.Config {
        cfg.EnableSplittingOnParsers = true
        return cfg
    }
}
```

Then, provide it in your Chaki app:

```go
app.Provide(fiberWrapper)
```

This allows customizing any Fiber config option.

## Usage

### Defining Controllers

Implement the `controller.Controller` interface:

```go
type YourController struct {
    *controller.Base
}

func NewYourController() controller.Controller {
    return &YourController{
        Base: controller.New("your-controller").SetPrefix("/api"),
    }
}

func (c *YourController) Routes() []route.Route {
    return []route.Route{
        route.Get("/endpoint", c.handler).Name("Endpoint Name"),
    }
}

func (c *YourController) handler(ctx context.Context, req YourRequest) (YourResponse, error) {
    // Handle request
    return response, nil
}
```

### Routes and Handlers

Use `route.Get`, `route.Post`, etc., with handler functions of type `func(context.Context, Req) (Res, error)`.

Request parsing is automatic for query, params, body, etc.

### Validation

Add validation tags to request structs (using `github.com/go-playground/validator/v10`).

Custom rules can be provided via `validation.Rule`.

### Request Parsing and Validation

Request structs can use tags to specify how fields are parsed from different parts of the HTTP request:

- `json:"field"` for body fields
- `query:"field"` for query parameters
- `param:"field"` for URL parameters
- `reqHeader:"Header-Name"` for request headers
- `cookie:"cookie_name"` for cookies
- `validate:"required"` or other validation rules

Example request struct:

```go
type ExampleRequest struct {
    Text           string `json:"text" validate:"required,min=3"`
    RepeatTimes    int    `json:"repeatTimes" validate:"required,gte=1"`
    NecessaryParam string `query:"necessaryParam" validate:"required"`
    CustomHeader   string `reqHeader:"X-Custom-Header" validate:"required"`
    SessionToken   string `cookie:"session_token" validate:"required"`
}
```

In your handler, the server module automatically parses these fields from the request. Validation is performed if tags are present.

### Responses

Handlers return data or `response.Responser` for custom responses.

### Advanced Middleware

Custom error handlers and logging middlewares.

#### Custom Error Handler Example

You can provide a custom error handler to the server module using `server.WithErrorHandler(yourHandler)`. This allows you to customize error responses, logging, and metadata.

Here's a generalized example of a custom error handler:

```go
package customerror

import (
	"errors"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

type CustomError struct {
	Timestamp     int64        `json:"timestamp"`
	CorrelationID string       `json:"correlationId"`
	Exception     string       `json:"exception"`
	RequestUri    string       `json:"requestUri"`
	RequestMethod string       `json:"requestMethod"`
	Status        int          `json:"status"`
	Errors        []ErrorDetail `json:"errors"`
	RawError      string       `json:"-"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func CustomErrHandler(c *fiber.Ctx, err error) error {
	fiberErr := &fiber.Error{}
	customErr := &CustomError{}

	if errors.As(err, &customErr) {
		setCustomErrorMetadata(c, customErr)
	} else if errors.As(err, &fiberErr) {
		customErr = fiberErrToCustomError(c, fiberErr)
	} else {
		customErr = unexpectedCustomError(c, err)
	}

	msg := resolveMessage(customErr)
	zapFields := []zap.Field{
		zap.String("rawError", customErr.RawError),
		zap.String("exception", customErr.Exception),
	}

	if msg == fiber.ErrUnauthorized.Message {
		logger.From(c.UserContext()).Warn(msg, zapFields...)
	} else {
		logger.From(c.UserContext()).Error(msg, zapFields...)
	}

	return c.Status(customErr.Status).JSON(customErr)
}

func setCustomErrorMetadata(c *fiber.Ctx, customError *CustomError) {
	customError.CorrelationID = ctxvaluer.CorrelationID.Get(c.UserContext(), "")
	customError.RequestUri = c.OriginalURL()
	customError.RequestMethod = c.Method()
}

func unexpectedCustomError(c *fiber.Ctx, err error) *CustomError {
	return &CustomError{
		Timestamp:     time.Now().UnixMilli(),
		CorrelationID: ctxvaluer.CorrelationID.Get(c.UserContext(), ""),
		Exception:     err.Error(),
		RequestUri:    c.OriginalURL(),
		RequestMethod: c.Method(),
		Status:        fiber.StatusInternalServerError,
		Errors: []ErrorDetail{
			{Message: "Unexpected error"},
		},
		RawError: err.Error(),
	}
}

func fiberErrToCustomError(c *fiber.Ctx, err *fiber.Error) *CustomError {
	return &CustomError{
		Timestamp:     time.Now().UnixMilli(),
		CorrelationID: ctxvaluer.CorrelationID.Get(c.UserContext(), ""),
		Exception:     err.Error(),
		RequestUri:    c.OriginalURL(),
		RequestMethod: c.Method(),
		Errors: []ErrorDetail{
			{Message: err.Error()},
		},
		Status: err.Code,
	}
}

func resolveMessage(customErr *CustomError) string {
	if len(customErr.Errors) > 0 && customErr.Errors[0].Message != "" {
		return customErr.Errors[0].Message
	}
	return customErr.Exception
}
```

To use this in your app:

```go
app.Use(server.Module(
  server.WithErrorHandler(customerror.CustomErrHandler),
))
```

This handler converts errors to a custom structure, adds metadata like correlation ID and request details, logs appropriately, and returns a JSON response.

#### Logging Middleware Example

You can create custom logging middlewares to add request-specific details to the logger context. This enhances logs with information like method, path, and custom headers.

Here's an example logging middleware:

```go
func LoggerMiddleware(c *fiber.Ctx) error {
	ctx := c.UserContext()
	ctx = logger.WithLogger(ctx, logger.From(ctx).With(
		zap.String("method", c.Method()),
		zap.String("request-path", c.Path()),
		zap.String("request-path-template", c.Route().Path),
		zap.String("request-uri", c.OriginalURL()),
		zap.String("platform", c.Get("platform")),
	))
	c.SetUserContext(ctx)
	return c.Next()
}
```

To use this middleware, add it via `server.WithMiddlewares(LoggerMiddleware)` or include it in a middleware group.

```go
app.Use(server.Module(
  server.WithMiddlewares(LoggerMiddleware),
))
```

This middleware enriches the context logger with request details for better traceability in logs.

### Response Customizations

Extended entity handling in response/entity.go.

#### Dynamic Response Status Example

Responses can implement a `Status() int` method to dynamically determine the HTTP status code based on the response data. This is useful for conditional status codes like 201 Created or 204 No Content.

Here's a generalized example:

```go
type ExampleResponse struct {
	Data string `json:"data"`
}

func NewNoContentResponse() *ExampleResponse {
	return &ExampleResponse{}
}

func (r *ExampleResponse) Status() int {
	if r.Data == "" {
		return fiber.StatusNoContent
	}
	return fiber.StatusCreated
}
```

In your handler, return this struct. The server module will call `Status()` to set the response code accordingly.

For example, in a controller handler:

```go
func (c *YourController) someHandler(ctx context.Context, req SomeRequest) (*ExampleResponse, error) {
	if /* no data condition */ {
		return NewNoContentResponse(), nil
	}
	return &ExampleResponse{Data: "Some value"}, nil
}
```

This allows flexible, data-driven response status without manual setting in each handler.

## Examples

### Basic Server

```go
package main

import (
    "github.com/Trendyol/chaki"
    "github.com/Trendyol/chaki/modules/server"
    "github.com/Trendyol/chaki/modules/server/controller"
    "github.com/Trendyol/chaki/modules/server/route"
)

func main() {
    app := chaki.New()
    app.Use(server.Module())
    app.Provide(NewHelloController)
    app.Start()
}

type HelloController struct { *controller.Base }

func NewHelloController() controller.Controller {
    return &HelloController{Base: controller.New("hello").SetPrefix("/hello")}
}

func (c *HelloController) Routes() []route.Route {
    return []route.Route{route.Get("/", c.hello)}
}

func (c *HelloController) hello(ctx context.Context, _ route.NoParam) (string, error) {
    return "Hello World", nil
}
```

### Advanced Controller with Fiber Routes

This example shows how to define a controller that uses raw Fiber handlers for routes, including GET, POST, and DELETE methods with request parsing.

```go
package controller

import (
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/route"
	"github.com/gofiber/fiber/v2"
)

type FilterRequest struct {
	// Fields for filtering
}

type FilterResponse struct {
	// Response data
}

type ActionRequest struct {
	// Fields for action
}

type Result struct {
	Success bool
}

func NewExampleController(exampleService ExampleService) *ExampleController {
	c := &ExampleController{
		controller.New("example").SetPrefix("/api"),
		exampleService,
	}
	return c
}

type ExampleController struct {
	*controller.Base
	exampleService ExampleService
}

func (ctr *ExampleController) Routes() []route.Route {
	return []route.Route{
		route.FiberGet[FilterRequest, *FilterResponse]("/filter", ctr.FilterHandler).Name("Filter"),
		route.FiberPost[ActionRequest, Result]("/:id1/:id2", ctr.ActionHandler).Name("Action"),
		route.FiberDelete[ActionRequest, Result]("/:id1/:id2", ctr.UndoActionHandler).Name("UndoAction"),
	}
}

func (c ExampleController) FilterHandler(ctx *fiber.Ctx) error {
	filter := &FilterRequest{}
	err := ctx.QueryParser(filter)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(c.exampleService.Filter(ctx.UserContext(), filter))
}

func (c ExampleController) ActionHandler(ctx *fiber.Ctx) error {
	id1 := ctx.Params("id1")
	id2 := ctx.Params("id2")

	result := c.exampleService.PerformAction(ctx.UserContext(), id1, id2)

	return ctx.JSON(&Result{Success: result})
}

func (c ExampleController) UndoActionHandler(ctx *fiber.Ctx) error {
	id1 := ctx.Params("id1")
	id2 := ctx.Params("id2")

	result := c.exampleService.UndoAction(ctx.UserContext(), id1, id2)

	return ctx.JSON(&Result{Success: result})
}
```

This demonstrates mixing Chaki's route definitions with direct Fiber context handling for more control over the request/response cycle.

### With Request Validation

Add tags like `validate:"required"` to request fields.

## Integrations

- **Swagger**: Add `swagger.Module()` to auto-generate API docs.
- **OTEL**: Use `otel.Module(otelserver.WithServer())` for tracing.
- **New Relic**: Use `newrelic.Module(nrserver.WithServer())` for monitoring.
- **ORM**: Controllers can inject ORM providers for database access.

## Troubleshooting

- **Validation not working**: Ensure custom rules are provided correctly and validator is initialized.
- **Panic recovery**: The module includes a recover middleware by default.
- **CORS issues**: Configure `cors` settings properly.
- **Logging**: Enable `server.logging` for request logs.
- **Health checks**: Use provided liveness and readiness endpoints for probes.
- Common pitfalls: Forgetting to set prefixes or names, improper request parsing for methods without body (e.g., GET).

For more details, check the source code or examples in the repository.
