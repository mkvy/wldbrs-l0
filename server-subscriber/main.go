package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/mkvy/wldbrs-l0/server-subscriber/database"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
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
	var orderData model.OrderData
	sub, err := sc.Subscribe(channel, func(msg *stan.Msg) {
		er := orderData.Scan(msg.Data)
		time.Sleep(time.Second)
		if er != nil {
			//TODO: catch the err if bad structure
			fmt.Println(er)
		}
		fmt.Printf("Received a message: %s\n", orderData)
	}, stan.StartWithLastReceived())
	defer sub.Unsubscribe()
	if err != nil {
		panic(err)
	}
  
	time.Sleep(time.Second * 3)
	validate := validator.New()
	err = validate.Struct(orderData)
	if err != nil {
		panic(err)
	}

	itemData := new(model.DataItem)
	itemData.OrderData = orderData
	itemData.ID = orderData.OrderUid
	fmt.Println("DATA", itemData)

	db, err := database.InitDBConn(db_driverName)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//_, err = db.SaveJsonToDB(itemData)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 1)

	fmt.Println("GETTING ALL ORDERS -----------------------------------")
	rows, err := db.GetAllOrders()
	defer rows.Close()
	strs := []model.DataItem{}
	for rows.Next() {
		str := model.DataItem{}
		err := rows.Scan(&str.ID, &str.OrderData)
		if err != nil {
			panic(err)
		}
		strs = append(strs, str)
	}
	for _, s := range strs {
		fmt.Println(s.ID, s.OrderData)
	}
	fmt.Println("----------- by id ------------")
	row := db.GetOrderByID("b563feb7b2b84b6test")
	rowData := new(model.DataItem)
	err = row.Scan(&rowData.ID, &rowData.OrderData)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowData.ID, rowData.OrderData)
}
