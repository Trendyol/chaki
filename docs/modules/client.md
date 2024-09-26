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
