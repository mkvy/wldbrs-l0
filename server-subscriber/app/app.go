package app

import (
	"context"
	"fmt"
	"github.com/mkvy/wldbrs-l0/server-subscriber/cache"
	"github.com/mkvy/wldbrs-l0/server-subscriber/config"
	"github.com/mkvy/wldbrs-l0/server-subscriber/database"
	"github.com/mkvy/wldbrs-l0/server-subscriber/server"
	"github.com/mkvy/wldbrs-l0/server-subscriber/store"
	"github.com/mkvy/wldbrs-l0/server-subscriber/subscriber"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	cacheService := cache.CacheInit()

	storeService := store.InitStore(*cacheService, *db)

	err = storeService.RestoreCache()

	if err != nil {
		log.Println("error restoring cache: db is empty")
	}

	sc := subscriber.CreateSub(*storeService)
	err = sc.Connect(app.cfg.Nats_server.Cluster_id, app.cfg.Nats_server.Client_id, app.cfg.Nats_server.Host+":"+app.cfg.Nats_server.Port)

	if err != nil {
		log.Println("Error while connecting to STAN")
	}

	sub, err := sc.SubscribeToChannel(app.cfg.Nats_server.Channel, stan.StartWithLastReceived())

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

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		defer server.Stop()
		defer sub.Unsubscribe()
		defer sc.Close()
		defer db.Close()
	}()

	if err := server.Srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
