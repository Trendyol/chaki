package main

import (
	"fmt"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/logger"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
		chaki.WithConfigReferencePath("secret", "secret.json"),
	)

	app.Provide()
	app.Invoke(printConfig)

	logger.Fatal(app.Start())
}

func printConfig(cfg *config.Config) {
	fmt.Println(
		"cfg var",
		cfg.GetString("val"),
		cfg.GetString("secret_val"),
	)
}
