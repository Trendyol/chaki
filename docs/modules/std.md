# Standard Modules

Standard modules are out-of-box modules required to develop applications in every aspect. These modules do not need to be registered using `app.Use` it comes with Chaki by default. And other modules like Kafka, server, and client, use std modules. These modules are main components for application development as well 

## Config Module

The config module reads the config file(s) using `viper` and provides a configuration instance for the injection pool. You can use multiple config files using config references, which can be optioned with a Chaki instance. Config module configures several formats without ease. Other modules can be configured using the Config module and config files so it's the main part of the Chaki and its config-first approach

### Usage
```go
//main.go
func main() {
    // create new instance
    app := chaki.New()

	app.WithOption( 
        // set config path
        // if config path not provided, default values is `resources/configs/config.yaml`
        chaki.WithConfigPath("config.yaml"),
        // set config reference path
        // config reference values are reachable in config.yaml using ${secret:<value.path>} annotations
		chaki.WithConfigReferencePath("secret", "secret.json"),
	)

    // invoke a basic function
	app.Invoke(printConfig)

    // start application
    if err := app.Start(); err != nil {
        logger.Fatal(err)
    }
}

func printConfig(cfg *config.Config) {
	fmt.Println(
		"cfg var",
		cfg.GetString("val"),
		cfg.GetString("secret_val"),
	)

    // Output: config var  example secret
}
```


```yaml
# config.yaml
val: exmaple
secret_val: ${secret:val}
```

```json
// secret.json
{
    "val": "secret"
}
```

## Logger Module

Logger module is built on top of the `uber-go/zap`. This module is automatically included as a standard module and reads configuration from the config file. It provides a logger instance that can be accessed via context and allows passing of default logging variables, which can be configured by other modules as well to provide extra information by default.

### Configuration

The logger module can be configured via the `logger` section in your config file:

```yaml
logger:
  timeKey: timestamp
  timeEncoder: iso8601
  level: info
  timezone: Europe/Istanbul
```

#### Configuration Options

- **timeKey** (string, default: `timestamp`): Specifies the key used for the timestamp in log output.

- **timeEncoder** (string, default: `epoch`): Specifies the time encoding format for log timestamps. Available options:
  - `epoch` - Unix epoch time (default)
  - `iso8601` - ISO8601 format
  - `rfc3339` - RFC3339 format
  - `rfc3339nano` - RFC3339 format with nanoseconds

- **level** (string, default: `info`): Sets the minimum log level. Available options:
  - `debug`
  - `info`
  - `warn`
  - `error`
  - `dpanic`
  - `panic`
  - `fatal`

- **timezone** (string, default: `Local`): Sets the timezone for log timestamps. Supports IANA location names, UTC, and UTC offsets. Available options:
  - `Local` - System's local timezone (default)
  - `UTC` - Universal Coordinated Time
  - `Europe/Istanbul`, `America/New_York`, etc. - IANA location names
  - `UTC+3`, `UTC-5`, etc. - UTC offsets in hours

### Context and Parameter Mapping

The Logger module, in conjunction with `ctxvaluer`, extracts HTTP headers into log fields. The `server.loggingHeaders` configuration uses a `field_name: Header-Name` format.

**Default Mappings:**
- `x-correlationId`: Correlation ID for request tracing.
- `x-executor-user`: Identity of the user performing the action.
- `x-agentname`: Name of the calling agent/service.
- `x-owner`: Owner of the resource or process.

These can be extended or overridden in `config.yaml`.

### Usage

The logger is automatically initialized when the application starts. You can access it from context:

```go
func foo(ctx context.Context) error {
	logger.From(ctx).Info("function executed")
	return nil
}
```