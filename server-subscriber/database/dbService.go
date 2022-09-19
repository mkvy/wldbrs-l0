package database

import (
	"database/sql"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
)

const (
	connStr = "user=root password=root dbname=testbd sslmode=disable"
)

type DBService struct {
	db *sql.DB
}

func InitDBConn(driverName string) (*DBService, error) {
	dbConn := DBService{}
	var err error
	dbConn.db, err = sql.Open(driverName, connStr)
	if err != nil {
		return &DBService{}, err
	}
	return &dbConn, err
}

func (dbService *DBService) Close() error {
	err := dbService.db.Close()
	return err
}

func (dbService *DBService) SaveJsonToDB(jsonData *model.DataItem) (sql.Result, error) {
	result, err := dbService.db.Exec(`insert into orders (id, orderdata) values ($1, $2)`, jsonData.ID, jsonData.OrderData)
	return result, err
}

func (dbService *DBService) GetAllOrders() ([]model.DataItem, error) {
	rows, err := dbService.db.Query("select * from orders")
	rowItem := model.DataItem{}
	rows.Scan(&rowItem.ID, &rowItem.OrderData)
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
	return strs, err
}

func (dbService *DBService) GetOrderByID(id string) (*model.DataItem, error) {
	row := dbService.db.QueryRow("select * from orders where id=$1", id)
	rowData := new(model.DataItem)
	err := row.Scan(&rowData.ID, &rowData.OrderData)
	return rowData, err
}
