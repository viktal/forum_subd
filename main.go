package main

import (
	"forum/application"
	//api "github.com/go-park-mail-ru/2020_2_MVVM.git/application"
	yconfig "github.com/rowdyroad/go-yaml-config"
)

func main() {
	var config api.Config
	yconfig.LoadConfig(&config, "configs/config.yaml", nil)
	app := api.NewApp(config)
	defer app.Close()
	app.Run()
}
