# App Composition

Use this page when wiring or refactoring Chaki application bootstrap.

## Recommended Flow

1. Create the application with `chaki.New()`.
2. Apply app-level options with `app.WithOption(...)` when config paths, config references, or timeouts need to change.
3. Register optional modules with `app.Use(...)`.
4. Register constructors and ready-made instances with `app.Provide(...)`.
5. Start the application with `app.Start()` and handle the returned error at the entrypoint.

## Bootstrap Example

```go
app := chaki.New()
app.WithOption(
	chaki.WithConfigPath("config/config.yaml"),
	chaki.WithConfigReferencePath("secret", "config/secret.json"),
	chaki.WithConfigReferencePath("kafka-auth", "config/kafka-auth.json"),
)

app.Use(
	server.Module(),
	client.Module(),
	kafka.Module(),
	swagger.Module(),
)

app.Provide(
	NewExampleClient,
	NewExampleService,
	NewExampleController,
	NewExampleConsumer,
)

if err := app.Start(); err != nil {
	panic(err)
}
```

- Keep each reference namespace aligned with the placeholders used in YAML, such as `${secret:...}` and `${kafka-auth:...}`.
- Real Chaki services often keep generic secrets and Kafka credentials in separate reference files.
- Use `app.Invoke(...)` only for startup hooks that do not fit normal constructor injection.

## Conventions

- Use `app.Use(...)` for framework modules such as server, client, Kafka, ORM, Couchbase, Swagger, OTEL, and New Relic.
- Use `app.Provide(...)` for repositories, services, controllers, consumers, and other constructors managed by dependency injection.
- Use `chaki.Valuer(...)` only when you already have a fully built instance that should be injected as-is.
- Keep the boot layer declarative. Module options and config should explain runtime behavior, not extra startup logic.

## Design Guidance

- Prefer small constructors that receive dependencies through Chaki injection.
- Enable only the modules and observability submodules the app actually needs.
- Keep repeated module keys, client names, consumer names, and similar identifiers in constants when reused across files.
- If a change needs secrets or externalized values, route them through config instead of embedding literals in code.
- If no config path override is provided, Chaki falls back to its default config path, so only add `WithConfigPath(...)` when the app really needs a non-default location.
