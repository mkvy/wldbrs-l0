package main

import (
	"ConnToNATSTest/publisher/service"
	"fmt"
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
		fmt.Println("Exiting")
		os.Exit(0)
	}
	fmt.Println(pathToPubData)
	nc := service.CreateSTAN()
	err := nc.Connect(clusterID, clientID, NATSStreamingURL)
	defer nc.Close()
	if err != nil {
		panic(err)
	}
	err = nc.PublishFromStdinCycle(channel)
	if err != nil {
		panic(err)
	}
}
