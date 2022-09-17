package database

import (
	"database/sql"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
)

const (
	connStr = "user=root password=root dbname=testbd sslmode=disable"
)

var id = 8

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

func (dbConn *DBConnection) SaveJsonToDB(jsonData model.OrderData) (sql.Result, error) {
	result, err := dbConn.db.Exec(`insert into foojsonb (foo_id, foo_note) values ($1, $2)`, id, jsonData)
	id++
	return result, err
}

func (dbConn *DBConnection) GetAllOrders() (*sql.Rows, error) {
	rows, err := dbConn.db.Query("select * from foojsonb")
	str := DBSchema{}
	rows.Scan(&str.Foo_id, &str.Foo_note)
	return rows, err
}

func (dbConn *DBConnection) GetOrderByID(id int) *sql.Row {
	row := dbConn.db.QueryRow("select * from foojsonb where foo_id=$1", id)
	return row
}

//getOrder cache restore here?
