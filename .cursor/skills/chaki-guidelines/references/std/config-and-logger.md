# Config And Logger

Use this page when the task needs configuration, secret references, or contextual logging.

## Config Rules

- Chaki's config module is built in and uses Viper-backed configuration.
- Use `chaki.WithConfigPath(...)` when the app should read a non-default config file.
- Use `chaki.WithConfigReferencePath(referenceName, path)` when config values should resolve `${referenceName:key}` lookups from another file.
- Inject `*config.Config` where code needs runtime configuration instead of opening config files manually.

## Config Reference Example

```yaml
serviceToken: ${secret:service_token}

client:
  example-client:
    baseurl: https://service.example.com
    timeout: 3s
```

```json
{
  "service_token": "replace-me"
}
```

```go
app := chaki.New()
app.WithOption(
	chaki.WithConfigPath("config/config.yaml"),
	chaki.WithConfigReferencePath("secret", "config/secret.json"),
)
```

- The reference name in `WithConfigReferencePath("secret", ...)` must match the `${secret:...}` prefix in YAML.
- Use separate reference namespaces when the app keeps generic secrets and Kafka credentials in different files.
- Use `cfg.Of("section")` for nested config groups and `config.ToStruct(...)` when a larger block should map into a struct.

## Logger Rules

- Use `logger.From(ctx)` for request-scoped or message-scoped logging.
- Treat logger enrichment as context enrichment. Middlewares and module hooks can add fields that later code should inherit.
- Keep log context close to boundaries such as HTTP middleware, Kafka consumers, and outbound clients.

## Design Guidance

- Prefer config-first module options over hardcoded transport, timeout, host, or credential values.
- If a config key or logger field name is reused across multiple files, keep it in a constant.
