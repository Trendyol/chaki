
<p align="center">
  <img src="./logo.png" width=256 />
</p>

![golang ci lint](https://github.com/Trendyol/chaki/actions/workflows/golangci-lint.yml/badge.svg) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Chaki

Chaki is an application builder framework designed to simplify the development process in Go by abstracting processes such as tracing, logging, Kafka integration, and Swagger documentation. It reduces the amount of boilerplate code, enhancing developer experience and productivity. Chaki handles complex tasks for you, allowing you to focus more on the core functionalities of your application.

Using Chaki, you can:

- Easily integrate the technologies and tools into your project with just a few configuration definitions.
- Quickly adapt to technological changes minumum effort on your stack, while continuing with your existing workflow.
- Create projects with code that primarily contains business logic, freeing you from technical bootstrap components.

## Usage

A simple usage example is provided below. For more detailed examples, please refer to the [examples](https://github.com/Trendyol/chaki) directory.

```go
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
