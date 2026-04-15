# Wrappers, Retry, And Errors

Use this page when the task affects client wrappers, custom error decoding, retries, or circuit breaker behavior.

## Error Handling

- Pass `client.WithErrDecoder(...)` to `factory.Get(...)` when the client needs custom response-to-error translation.
- Use `client.GenericClientError` when code needs structured access to client error details.
- Keep error mapping close to the client boundary so services receive domain-appropriate failures.

## Wrapper Rules

- Use `client.DriverWrapper` when requests need extra headers or driver-level behavior.
- Wrapper logic should read values from request context and project them into headers or Resty settings.
- Keep reused header names and context keys in constants.

## Retry And Circuit Rules

- Retry behavior is config-driven and automatically applied when enabled.
- Circuit breaker activation is config-driven, and requests must carry a command name when that circuit configuration is enabled.
- `RequestWithCommand(...)` is the normal way to attach that command name for a request that should run through the configured circuit breaker path.
- Keep retry presets, circuit presets, and thresholds in config rather than embedding them in request code.
- Avoid duplicate error logging. Chaki's default client error handling already avoids noisy double-logging.

## Retry And Circuit Example

```yaml
client:
  circuitPresets:
    - name: "Base Circuit"
      timeout: 1500
      maxConcurrentRequests: 1000
      errorPercentThreshold: 50
      requestVolumeThreshold: 6
      sleepWindow: 10000

  retryPresets:
    - name: "Fast Retry"
      count: 3
      interval: 100ms
      maxDelay: 2s
      delayType: constant

  inventory-api:
    baseurl: https://service.example.com
    circuit:
      enabled: true
      preset: "Base Circuit"
    retry:
      enabled: true
      preset: "Fast Retry"
```

```yaml
client:
  inventory-api:
    baseurl: https://service.example.com
    circuit:
      enabled: true
      preset: custom
      timeout: 3000
      maxConcurrentRequests: 10
      errorPercentThreshold: 25
      requestVolumeThreshold: 5
      sleepWindow: 5000
      statusCodeConfig:
        treatAllErrorCodesAsFailure: true
        ignoreStatusCodes: [404]
```

- Use `client.circuitPresets` and `client.retryPresets` for reusable presets across multiple clients.
- When `circuit.preset` is `custom`, the inline `circuit.*` fields are used instead of a named preset.
- Enabling circuit behavior in config is not enough by itself; call `RequestWithCommand(...)` so the request carries the circuit command name.
