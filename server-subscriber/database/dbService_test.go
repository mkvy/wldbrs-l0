package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
	"os"
	"testing"
)

func TestDBService_SaveOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	text, err := os.ReadFile("../testdata/payload.json")
	if err != nil {
		t.Fatal("error reading testdata file")
	}
	od := new(model.OrderData)
	err = od.Scan(text)
	itemData := new(model.DataItem)
	itemData.OrderData = *od
	itemData.ID = od.OrderUid
	mock.ExpectExec("insert into orders").WithArgs(od.OrderUid, od).WillReturnResult(sqlmock.NewResult(0, 1))
	dbS := NewDB(db)
	_, err = dbS.SaveOrder(itemData)
	if err != nil {
		t.Errorf("there were error saving data: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBService_Close(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectClose()
	dbs := NewDB(db)
	dbs.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
