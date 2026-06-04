# Generation And Access Control

Use this page when Swagger output or Swagger visibility needs to change.

## Wiring Pattern

- Register `swagger.Module()` alongside `server.Module()`.
- Rely on route metadata, request tags, and response types to drive generated documentation.

## Generation Rules

- Query, path, and body tags directly affect generated schema and parameters.
- If generated docs look wrong, inspect the route request model and metadata before adding manual overrides.

## Access Control Rules

- Configure host restrictions under `server.swagger.blockedHostsContains` and `server.swagger.allowedHostsContains`.
- These rules are substring checks against the effective host, not exact hostname equality checks.
- The blocked list takes precedence over the allowed list.
- Denied access returns `404` to avoid exposing the existence of Swagger endpoints.
- Consider forwarded host behavior when the app sits behind a reverse proxy or API gateway.

## Config Example

```yaml
server:
  swagger:
    blockedHostsContains:
      - "api.example.com"
      - "external.example.com"
    allowedHostsContains:
      - "internal.example.com"
      - "localhost"
      - "127.0.0.1"

info:
  description: "Example Application"
  version: "1.0.0"
  title: Example Application
  terms: https://swagger.io/terms/
  email: "team@example.com"
```

- Put Swagger host restrictions under the `server.swagger` branch, not at the root of the config file.
- Use the `info` block to shape generated API documentation metadata.
- If both host lists are present, the blocked list wins.
- Because matching uses substring checks and may read `X-Forwarded-Host`, keep the allowed and blocked values specific enough for the deployment topology.
