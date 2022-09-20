package main

import (
	_ "github.com/lib/pq"
	app "github.com/mkvy/wldbrs-l0/server-subscriber/app"
	config "github.com/mkvy/wldbrs-l0/server-subscriber/config"
)

const (
	NATSStreamingURL = "127.0.0.1:4223"
	clusterID        = "test-cluster"
	clientID         = "test-client"
	channel          = "testch"
	db_driverName    = "postgres"
	addr_server      = "localhost:8181"
)

func main() {
	config := new(config.Config)
	config.InitFile()
	app := app.InitApp(*config)
	app.Run()
}
