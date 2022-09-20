package app

import (
	"fmt"
	"github.com/mkvy/wldbrs-l0/server-subscriber/cache"
	"github.com/mkvy/wldbrs-l0/server-subscriber/config"
	"github.com/mkvy/wldbrs-l0/server-subscriber/database"
	"github.com/mkvy/wldbrs-l0/server-subscriber/server"
	"github.com/mkvy/wldbrs-l0/server-subscriber/store"
	"github.com/mkvy/wldbrs-l0/server-subscriber/subscriber"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

type App struct {
	cfg config.Config
}

func InitApp(cfg config.Config) *App {
	app := App{}
	app.cfg = cfg
	return &app
}

func (app *App) Run() {
	db, err := database.InitDBConn(app.cfg)
	if err != nil {
		log.Println("Error while connnecting to database")
		panic(err)
	}
	defer db.Close()

	cacheService := cache.CacheInit()

	storeService := store.InitStore(*cacheService, *db)

	err = storeService.RestoreCache()

	if err != nil {
		log.Println("error restoring cache: db is empty")
	}

	sc := subscriber.CreateSub(*storeService)
	err = sc.Connect(app.cfg.Nats_server.Cluster_id, app.cfg.Nats_server.Client_id, app.cfg.Nats_server.Host+":"+app.cfg.Nats_server.Port)
	defer sc.Close()
	if err != nil {
		log.Println("Error while connecting to STAN")
	}

	sub, err := sc.SubscribeToChannel(app.cfg.Nats_server.Channel, stan.StartWithLastReceived())

	defer sub.Unsubscribe()
	if err != nil {
		log.Println("Error while subscribing to channel")
	}

	time.Sleep(time.Second * 3)

	fmt.Println("getting from cache ", storeService.GetFromCacheByUID("b563feb7b2b84b6test"))

	dItems, err := storeService.GetAllOrders()
	if err != nil {
		log.Println("orders not found in database")
		log.Println(err)
	}

	fmt.Println("get all orders from database ", dItems)

	server := server.InitServer(*storeService, app.cfg.Http_server.Host+":"+app.cfg.Http_server.Port)
	err = server.Start()
	if err != nil {
		log.Println("error while starting server")
		panic(err)
	}
	defer server.Stop()
}
