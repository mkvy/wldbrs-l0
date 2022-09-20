package main

import (
	_ "github.com/lib/pq"
	app "github.com/mkvy/wldbrs-l0/server-subscriber/app"
	config "github.com/mkvy/wldbrs-l0/server-subscriber/config"
)

func main() {
	config := new(config.Config)
	config.InitFile()
	app := app.InitApp(*config)
	app.Run()
}
