package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mkvy/wldbrs-l0/server-subscriber/database"
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
	var data string
	sub, err := sc.Subscribe(channel, func(msg *stan.Msg) {
		data = string(msg.Data)
		fmt.Printf("Received a message: %s\n", data)
	}, stan.StartWithLastReceived())
	defer sub.Unsubscribe()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 8)

	db, err := database.InitDBConn(db_driverName)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 1)
	fmt.Println("GETTING ALL ORDERS -----------------------------------")
	rows, err := db.GetAllOrders()
	defer rows.Close()
	strs := []database.DBSchema{}
	for rows.Next() {
		str := database.DBSchema{}
		err := rows.Scan(&str.Foo_id, &str.Foo_note)
		if err != nil {
			panic(err)
		}
		strs = append(strs, str)
	}
	for _, s := range strs {
		fmt.Println(s.Foo_id, s.Foo_note)
	}
	fmt.Println("----------- by id ------------")
	row := db.GetOrderByID(5)
	rowData := database.DBSchema{}
	err = row.Scan(&rowData.Foo_id, &rowData.Foo_note)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowData.Foo_id, rowData.Foo_note)
}
