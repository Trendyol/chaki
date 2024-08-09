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

Logger module is built on top of the `uber-go/zap`. This module can start a logger using context and allows to passing of default logging variables, these logging variables can be configured by other modules as well to provide extra information by default

```go
func foo(ctx context.Context) error {
	logger.From(ctx).Info("function executed")
}
```