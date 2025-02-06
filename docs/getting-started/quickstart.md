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

Chaki is built on the `uber-go/fx` dependency injection system and extends its functionality with a robust module system. This simple yet powerful approach allows developers to enhance their applications with ease by leveraging pre-built modules.

### Key Features of Modules

- **Simple Integration**: Add modules effortlessly to your system using the `app.Use` method.
- **Direct Access**: Once a module is added, you can directly utilize the instances it provides—such as `*gorm.DB` or `sarama.SyncProducer`—in any constructor within your application.
- **Flexible Configuration**: Customize modules easily through configuration files, ensuring seamless adaptation to your application's needs.
