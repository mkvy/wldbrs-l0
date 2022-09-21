package database

import (
	"database/sql"
	"github.com/mkvy/wldbrs-l0/server-subscriber/config"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
)

type DBService struct {
	db *sql.DB
}

func NewDB(database *sql.DB) *DBService {
	return &DBService{db: database}
}

func InitDBConn(cfg config.Config) (*DBService, error) {
	dbConn := DBService{}
	var err error
	connStr := "user=" + cfg.Database.Username + " password=" + cfg.Database.Password + " dbname=" + cfg.Database.DBname + " sslmode=disable"
	dbConn.db, err = sql.Open(cfg.Database.DriverName, connStr)
	if err != nil {
		return &DBService{}, err
	}
	return &dbConn, err
}

func (dbService *DBService) Close() error {
	err := dbService.db.Close()
	return err
}

func (dbService *DBService) SaveOrder(jsonData *model.DataItem) (sql.Result, error) {
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
			return strs, err
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
