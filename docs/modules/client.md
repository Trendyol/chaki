# Client

The client module enables your application to interact with other web services over HTTP using the [Resty](https://github.com/go-resty/resty) library.

This module simplifies the process of creating an HTTP client. It reads the base URL and timeout configurations from a config file. To create a client, you need to add the module to your application as follows:

```go

    app := chaki.New()

	app.Use(
        // ...
		client.Module(),
        // ...
	)

    //...
```

To create a client, you can use the following code:

```go
type exampleClient struct {
	*client.Base
}

func newClient(f *client.Factory) *exampleClient {
	return &exampleClient{
		Base: f.Get("example-client"),
	}
}
```

the `example-client` name should match the name in the config file to configure the client from the config file:

```yaml
client:
  example-client:
    baseUrl: "http://super-duper-client-url.com"
    timeout: 500ms
```

To create a request, you should use the following code. This ensures that tracing and other metadata are used on the request as all metadata is under context.

```go
func (cl *exampleClient) SendHello(ctx context.Context) (string, error) {
	resp := &response.Response[string]{}

    request := cl.Request(ctx) // this gives you an *resty.Request, to work with.

	if _, err := request.
        SetResult(resp).
        Get("/hello"); err != nil {
		return "", err
	}

	return resp.Data, nil
}
```

If you want to log every outgoing request and incoming response, you can simply set `logging` key to `true` on config.

```yaml
client:
example-client:
  baseUrl: "http://super-duper-client-url.com"
  timeout: 500ms
  logging: true
```

---

## Error Handler

By default, Chaki provides a built-in error handler to encapsulate incoming errors. The source code can be found in `modules/client/errors.go`. To avoid log chaos, error cases are not logged by default.

To access the details of the errors, you can cast the error type into `GenericClientError` type as follows:

```go

    _, err := cl.SendSomeRequest()
    genericError := client.GenericClientError{}
    errors.As(err, genericError)
    logger.From(ctx).Error(genericError.ClientName)

```

### Providing error handler

You can provide a custom error handler to handle errors in a more specific way. The error handler function should accept a `context.Context` and a `*resty.Response` as parameters.
But returning an error that implements `Statuser` from the server module _having a *Status() int* method will help you to return correct status code from your endpoint._
```go
func newClient(f *client.Factory) *exampleClient {
	return &exampleClient{
		Base: f.Get("example-client", client.WithErrDecoder(customErrorDecoder)),
	}
}

func customErrorDecoder(_ context.Context, res *resty.Response) error {
	if res.StatusCode() == 404 {
		return fmt.Errorf("not found")
	}
	return nil
}
```

---

## Wrappers

You can add wrappers to clients to extend their functionality. Chaki provides a default wrapper that adds the following headers to requests if the corresponding values are present in the context:

```go
    CorrelationIDKey = "x-correlationId"
	ExecutorUserKey  = "x-executor-user"
	AgentNameKey     = "x-agentname"
	OwnerKey         = "x-owner"
```

### Providing an wrapper

You can wrap the existing client as follows.

```go


type user struct {
	publicUsername string
	publicTag      string
}

func HeaderWrapper() client.DriverWrapper {
	return func(restyClient *resty.Client) *resty.Client {
		return restyClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			ctx := r.Context()

			h := map[string]string{}

			if v := ctx.Value("user"); v != nil {
				user := v.(user)
				h["publicUsername"] = user.publicUsername
				h["publicTag"] = user.publicTag
			}

			r.SetHeaders(h)
			return nil
		})
	}
}

func newClient(f *client.Factory) *exampleClient {
	return &exampleClient{
		Base: f.Get("example-client",
			client.WithDriverWrappers(HeaderWrapper())),
	}
}

```

## Circuit Breaker

**-Currently WIP-**

The client module includes a built-in circuit breaker functionality using Hystrix-go with predefined circuit presets and ability to add some custom settings.

This feature is turned-off by default. To enable it, you can use the following configurations.

```yaml
client:
  circuitPresets:
    - name: "Custom Preset - 1"
      timeout: 1000 # in ms
      maxConcurrentRequests: 1250
      errorPercentThreshold: 10
      requestVolumeThreshold: 25
      sleepWindow: 5000 # in ms
  predefinedPresetClient:
    baseUrl: "http://super-duper-client-url.com"
    circuit:
      enabled: true
      preset: "default" # Available presets: default, aggressive, relaxed
  customCircuitClient:
    baseUrl: "http://super-duper-client-url.com"
    circuit:
      enabled: true
      preset: "Custom Preset - 1"
  inlineCustomPreset:
    baseUrl: "http://super-duper-client-url.com"
    circuit:
      preset: custom # it needs to be 'custom' spesifically
      enabled: true
      timeout: 2000
      maxConcurrentRequests: 50
      errorPercentThreshold: 10
      requestVolumeThreshold: 20
      sleepWindow: 1000
```

Even if you configure circuit breaker settings in your configuration files, you need to explicitly specify which requests should be protected by the circuit breaker by using `RequestWithCommand` instead of regular `Request`.
Commands are scoped to their respective clients. This means that even if you use the same command name across different clients, they are treated as separate and independent circuit breakers.

```go
// This won't use circuit breaker even if circuit breaker is configured
client.Request(ctx).Get("/api/users")

// This will use circuit breaker with command name "get-users"
client.RequestWithCommand(ctx, "get-users").Get("/api/users")
```

### Built-in presets

- **default**: Moderate settings (5s timeout, 100 concurrent requests)
- **aggressive**: Strict settings (2s timeout, 50 concurrent requests)
- **relaxed**: Lenient settings (10s timeout, 200 concurrent requests)

## Retry

Even if you configure retry settings in your configuration files, all requests will automatically use the configured retry mechanism. Unlike circuit breaker, you don't need to specify any special method - the retry mechanism works automatically based on your configuration.

### Configuration:

```yaml
client:
  service-name:
    retry:
      enabled: true
      preset: "default" # Available presets: default, exponential, aggressive, aggressiveExponential, relaxed, relaxedExponential
```

### Custom Configuration:

```yaml
client:
  retryPresets:
    - name: "Custom Preset - 1"
      count: 3
      interval: 100ms
      maxDelay: 2s
      delayType: constant
  customRetryClient:
    retry:
      enabled: true
      preset: "custom"
      count: 3 # Number of retry attempts
      interval: "100ms" # Base interval between retries
      maxDelay: "5s" # Maximum delay cap for exponential backoff
      delayType: "constant" # or "exponential"
  customPresetClient:
    retry:
      enabled: true
      preset: "Custom Preset - 1"
```

### Delay Types:

1. **Constant Delay**:

   - Fixed time interval between retries
   - Example: 100ms -> 100ms -> 100ms

2. **Exponential Delay**:
   - Increases exponentially with each retry attempt
   - Includes jitter to prevent thundering herd
   - Example: 100ms -> 200ms -> 400ms (plus random jitter)
   - Capped by maxDelay setting

**_Unlike circuit breaker which requires explicit command specification, retry mechanism works automatically for all requests based on the client's configuration._**

### Built-in Presets:

1. **default**:

   - Count: 3 retries
   - Interval: 100ms
   - MaxDelay: 5s
   - DelayType: constant

2. **exponential**:

   - Count: 3 retries
   - Interval: 100ms
   - MaxDelay: 5s
   - DelayType: exponential

3. **aggressive**:

   - Count: 7 retries
   - Interval: 50ms
   - MaxDelay: 2s
   - DelayType: constant

4. **aggressiveExponential**:

   - Count: 7 retries
   - Interval: 50ms
   - MaxDelay: 2s
   - DelayType: exponential

5. **relaxed**:

   - Count: 2 retries
   - Interval: 500ms
   - MaxDelay: 2s
   - DelayType: constant

6. **relaxedExponential**:
   - Count: 2 retries
   - Interval: 500ms
   - MaxDelay: 2s
   - DelayType: exponential
