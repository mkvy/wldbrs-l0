package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"time"
)

const (
	NATSStreamingURL = "127.0.0.1:4223"
	clusterID        = "test-cluster"
	clientID         = "test-client"
	channel          = "testch"
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
	fmt.Println("Here is a data", data)

	connStr := "user=root password=root dbname=testbd sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	/*result, err := db.Exec(`insert into foo (foo_id, foo_note) values ('2', $1)`, data)
	if err != nil {
		panic(err)
	}*/
	//fmt.Println(result.RowsAffected())
	result, err := db.Exec(`insert into foojson (foo_id, foo_note) values ('3', $1)`, data)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	result, err = db.Exec(`insert into footext (foo_id, foo_note) values ('3', $1)`, data)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	result, err = db.Exec(`insert into foojsonb (foo_id, foo_note) values ('3', $1)`, data)
	if err != nil {
		panic(err)
	}
}
