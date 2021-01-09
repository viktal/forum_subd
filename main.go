package main

import (
	"forum/application"
	yconfig "github.com/rowdyroad/go-yaml-config"
)

func main() {
	var config api.Config
	yconfig.LoadConfig(&config, "configs/config.yaml", nil)
	app := api.NewApp(config)
	app.Run()
}
