package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mkvy/wldbrs-l0/server-subscriber/cache"
	"github.com/mkvy/wldbrs-l0/server-subscriber/database"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
	"github.com/mkvy/wldbrs-l0/server-subscriber/store"
	"github.com/nats-io/stan.go"
	"time"
)

const (
	NATSStreamingURL = "127.0.0.1:4223"
	clusterID        = "test-cluster"
	clientID         = "test-client"
	channel          = "testch"
	db_driverName    = "postgres"
)

func main() {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(NATSStreamingURL))
	defer sc.Close()
	if err != nil {
		panic(err)
	}

	db, err := database.InitDBConn(db_driverName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cacheService := cache.CacheInit()

	storeService := store.InitStore(*cacheService, *db)

	var orderData model.OrderData

	sub, err := sc.Subscribe(channel, func(msg *stan.Msg) {
		er := orderData.Scan(msg.Data)
		time.Sleep(time.Second)
		if er != nil {
			//TODO: catch the err if bad structure
			fmt.Println(er)
		}
		fmt.Printf("Received a message: %s\n", orderData)
		err = storeService.SaveOrderData(msg.Data)
		if err != nil {
			panic(err)
		}

	}, stan.StartWithLastReceived())
	defer sub.Unsubscribe()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 3)

	fmt.Println("getting from cache ", storeService.GetFromCacheByUID("b563feb7b2b84b6test"))

	dItems, err := storeService.GetAllOrders()
	if err != nil {
		panic(err)
	}

	fmt.Println("get all orders from database ", dItems)

}
