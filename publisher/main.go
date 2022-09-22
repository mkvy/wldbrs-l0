package main

import (
	"github.com/mkvy/wldbrs-l0/publisher/service"
	"log"
)

const (
	NATSStreamingURL = "localhost:4223"
	clusterID        = "test-cluster"
	clientID         = "test-publisher"
	channel          = "testch"
)

func main() {
	nc := service.CreateSTAN()
	err := nc.Connect(clusterID, clientID, NATSStreamingURL)
	defer nc.Close()
	if err != nil {
		log.Println("Error while connecting to nats")
		panic(err)
	}
	_ = nc.PublishFromStdinCycle(channel)
}
