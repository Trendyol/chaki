# Factory And Config

Use this page when creating a Chaki HTTP client or updating its config.

## Factory Pattern

- Inject `*client.Factory` into the constructor.
- Create the module-specific client with `factory.Get(clientName)`.
- Embed `*client.Base` in the client type if the client should expose Chaki request helpers.
- Keep `clientName` in a constant when it is reused in code and config.

## Client Config Example

```yaml
client:
  inventory-api:
    baseurl: https://service.example.com
    timeout: 3s
    logging: false
```

```go
const inventoryClientName = "inventory-api"

func NewInventoryClient(factory *client.Factory) *InventoryClient {
	return &InventoryClient{
		Base: factory.Get(inventoryClientName),
	}
}
```

- The YAML key under `client.<name>` must match the name passed to `factory.Get(...)`.
- Chaki's client driver reads the canonical key `baseurl`. Some service repositories use mixed casing such as `baseUrl`, but the skill should document `baseurl`.

## Config Rules

- Configure clients under `client.<name>`.
- Common built-in keys include `baseurl`, `timeout`, optional `logging`, and optional `debug`.
- Keep host, timeout, retry, and circuit settings in config instead of inline code.

## Request Rules

- Build requests with `Base.Request(ctx)` so tracing and contextual metadata move with the call.
- Use typed response holders with `SetResult(...)` where the response contract is known.
- Use `RequestWithCommand(ctx, commandName)` when circuit breaker command naming matters, and keep reused command names in constants.
