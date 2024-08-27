# Swagger Module

Swagger module is for generating swagger endpoints for route definitions. 

The Swagger module automatically generates Swagger endpoints for your route definitions, providing an interactive API documentation interface. This module relies on the `EndpointDefinition` conversion of the route definitions created by the Server module. Therefore, the Server module is required for the Swagger module to function.
The swagger documents are embedded in the application itself, there is no redirecting.

### Usage
To include the Swagger module in your Chaki application, add it to the `app.Use` function alongside the Server module:
```go
...

app := chaki.New()

app.Use(
    server.Module(),
    swagger.Module(),
)

....
```

### How it works?
The Server module generates and stores route definitions, each containing a Meta model, which includes metadata about the route.
```go
type Meta struct {
	Method        string
	Path          string
	Name          string
	Desc          string
	Func          fiber.Handler
	Req           reflect.Type
	Res           reflect.Type
	Middlewares   []fiber.Handler
	DefaultStatus int
}
```

The Swagger module converts this metadata into Swagger definitions and paths based on the request model. If a request model includes a JSON body parameter, the Swagger module creates a corresponding definition and adds a reference to it in the endpoint's parameters field. For query and path variables, it automatically includes them in the parameters field.


For example, given the following request model:
```go
type GreetWithBodyRequest struct {
    Text           string `json:"text" validate:""`
    RepeatTimes    int    `json:"repeatTimes" validate:"required"`
    NecessaryParam string `query:"necessaryParam" validate:"required"`
}
```

The Swagger module generates the following JSON:
```json lines
{
  "definitions": {
    "GreetWithBodyRequest": {
      "properties": {
        "repeatTimes": {
          "format": "int",
          "type": "integer"
        },
        "text": {
          "type": "string"
        }
      },
      "required": [
        "repeatTimes"
      ],
      "type": "object"
    }
  },
  
  /// ...

  "paths": {
    "/hello/body": {
      "post": {
        /// ...
        "parameters": [
          {
            "description": "",
            "in": "query",
            "name": "necessaryParam",
            "required": true,
            "type": "string"
          },
          {
            "description": "",
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/GreetWithBodyRequest"
            }
          }
        ],
        /// ...
    },}
    /// ...
  }
}

```

### Config
The Swagger module offers simple configurations to customize the generated API documentation. Below are the default values:

```yaml

info:
  description: ""
  version: 1.0.0
  title: Application
  terms: http://swagger.io/terms/
  email: ""
```
- Description: A brief description of your API. This field is optional.
- Version: The version of your API. Defaults to 1.0.0.
- Title: The title of your API documentation. Defaults to Application.
- Terms: A URL to the terms of service for your API. Defaults to http://swagger.io/terms/.
- Email: The contact email for API support or inquiries.




