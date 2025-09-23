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
The Swagger module offers configurations to customize the generated API documentation and control access. Below are the available configurations:

#### API Documentation Configuration
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

#### Host-based Access Control
To protect your Swagger documentation from unauthorized access, you can configure host-based restrictions. This is particularly useful when your application is exposed to both internal and external networks.

```yaml
server:
  swagger:
    # Block access when hostname contains any of these substrings
    blockedHostsContains:
      - "api.example.com"
      - "external.company.com"
    
    # Optional: Only allow access when hostname contains any of these substrings
    # If both allowedHostsContains and blockedHostsContains are set,
    # blockedHostsContains takes precedence for better security
    allowedHostsContains:
      - "internal.company.com"
      - "localhost"
      - "127.0.0.1"
```

**Configuration Options:**
- `blockedHostsContains`: Array of strings. If the request hostname contains any of these substrings, access to Swagger will be denied with a 404 response. Takes precedence over allowlist for security.
- `allowedHostsContains`: Array of strings. If specified, only requests from hostnames containing these substrings will be allowed, unless they are also in the blocked list.

**Host Detection:**
The middleware checks both the `Host` header and `X-Forwarded-Host` header (common when behind reverse proxies). The `X-Forwarded-Host` header takes precedence if present.

**Example Use Cases:**

1. **Block External Gateway Access:**
```yaml
server:
  swagger:
    blockedHostsContains:
      - "apigw.company.com"
      - "public.api.company.com"
```

2. **Internal-Only Access:**
```yaml
server:
  swagger:
    allowedHostsContains:
      - "internal"
      - "corp"
      - "localhost"
      - "127.0.0.1"
```

3. **Development Environment:**
```yaml
server:
  swagger:
    allowedHostsContains:
      - "localhost"
      - "127.0.0.1"
      - "dev.company.com"
      - "staging.company.com"
```

**Security Note:** When access is denied, the middleware returns a 404 status code instead of 403 to avoid revealing the existence of Swagger endpoints to unauthorized users.




