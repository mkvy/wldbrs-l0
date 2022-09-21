package main

import (
	"github.com/mkvy/wldbrs-l0/publisher/service"
	"log"
	"os"
)

const (
	NATSStreamingURL = "127.0.0.1:4223"
	clusterID        = "test-cluster"
	clientID         = "test-publisher"
	channel          = "testch"
)

func main() {
	var pathToPubData string
	if len(os.Args) > 0 {
		pathToPubData = os.Args[1]
	} else {
		log.Println("Exiting")
		os.Exit(0)
	}
	log.Println(pathToPubData)
	nc := service.CreateSTAN()
	err := nc.Connect(clusterID, clientID, NATSStreamingURL)
	defer nc.Close()
	if err != nil {
		log.Println("Error while connecting to nats")
		panic(err)
	}
	//err = nc.PublishFromFile(channel, pathToPubData)
	_ = nc.PublishFromStdinCycle(channel)
}
