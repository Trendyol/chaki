## Install
```
go get github.com/Trendyol/chaki
```

## Usage

A simple usage example is provided below. For more detailed examples, please refer to the [examples](https://github.com/Trendyol/chaki/tree/main/example) directory.

```golang
func main() {
	// create instance
	app := chaki.New()

	// add modules
	app.Use(
	 foomodule.Module(),
	)

	// provide constructors
	app.Provide(
	  NewFooRepository,
	  NewFooService,
	  NewBarClient,
	  NewFooController,
	  NewFooCreatedConsumer,
	  // or you can just provide an already created instance
	  chaki.Valuer(barRepositoryInstance),
	)

	// start application
	if err := app.Start(); err != nil {
	    logger.Fatal(err)
	}
}
```


## Modules

Chaki is fundamentally based on the `uber-go/fx` dependency injection system and includes additional features such as a module system. This framework, built on a simple foundation, derives much of its power from its modules.

- Modules are easily added to the system with `app.Use`.
- Once added, you can directly use instances produced by the module, such as `*gorm.Gorm` or `sarama.SyncProducer`, in any constructor you want.
- Modules can be easily configured using configuration files.