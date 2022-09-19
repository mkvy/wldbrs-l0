package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mkvy/wldbrs-l0/server-subscriber/cache"
	"github.com/mkvy/wldbrs-l0/server-subscriber/database"
	"github.com/mkvy/wldbrs-l0/server-subscriber/server"
	"github.com/mkvy/wldbrs-l0/server-subscriber/store"
	"time"
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

	db, err := database.InitDBConn(db_driverName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cacheService := cache.CacheInit()

	storeService := store.InitStore(*cacheService, *db)

	err = storeService.RestoreCache()
	if err != nil {
		panic(err)
	}

	/*sc := subscriber.CreateSub(*storeService)
	err = sc.Connect(clusterID, clientID, NATSStreamingURL)
	defer sc.Close()
	if err != nil {
		panic(err)
	}

	sub, err := sc.SubscribeToChannel(channel, stan.StartWithLastReceived())

	defer sub.Unsubscribe()
	if err != nil {
		panic(err)
	}
	*/
	time.Sleep(time.Second * 3)

	fmt.Println("getting from cache ", storeService.GetFromCacheByUID("b563feb7b2b84b6test"))

	dItems, err := storeService.GetAllOrders()
	if err != nil {
		panic(err)
	}

	fmt.Println("get all orders from database ", dItems)

	server := server.InitServer(*storeService, addr_server)
	err = server.Start()
	if err != nil {
		panic(err)
	}
	defer server.Stop()
}
