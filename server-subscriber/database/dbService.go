package database

import (
	"database/sql"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
)

const (
	connStr = "user=root password=root dbname=testbd sslmode=disable"
)

type DBSchema struct {
	Foo_id   int
	Foo_note string
}

type DBConnection struct {
	db *sql.DB
}

func InitDBConn(driverName string) (*DBConnection, error) {
	dbConn := DBConnection{}
	var err error
	dbConn.db, err = sql.Open(driverName, connStr)
	if err != nil {
		return &DBConnection{}, err
	}
	return &dbConn, err
}

func (dbConn *DBConnection) Close() error {
	err := dbConn.db.Close()
	return err
}

func (dbConn *DBConnection) SaveJsonToDB(jsonData *model.DataItem) (sql.Result, error) {
	result, err := dbConn.db.Exec(`insert into orders (id, orderdata) values ($1, $2)`, jsonData.ID, jsonData.OrderData)
	return result, err
}

func (dbConn *DBConnection) GetAllOrders() (*sql.Rows, error) {
	rows, err := dbConn.db.Query("select * from orders")
	rowItem := model.DataItem{}
	rows.Scan(&rowItem.ID, &rowItem.OrderData)
	return rows, err
}

func (dbConn *DBConnection) GetOrderByID(id string) *sql.Row {
	row := dbConn.db.QueryRow("select * from orders where id=$1", id)
	return row
}

//getOrder cache restore here?
