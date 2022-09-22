package store

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mkvy/wldbrs-l0/server-subscriber/cache"
	"github.com/mkvy/wldbrs-l0/server-subscriber/database"
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
	"log"
)

type StoreService struct {
	cache cache.CacheService
	db    database.DBService
}

func InitStore(cache cache.CacheService, db database.DBService) *StoreService {
	StoreService := StoreService{
		cache: cache,
		db:    db,
	}
	return &StoreService
}

func (ss *StoreService) SaveOrderData(data []byte) error {
	od := new(model.OrderData)
	err := od.Scan(data)
	if err != nil {
		log.Println("Wrong format")
		return err
	}
	validate := validator.New()
	err = validate.Struct(od)
	if err != nil {
		log.Println(err)
		return err
	}
	itemData := new(model.DataItem)
	itemData.OrderData = *od
	itemData.ID = od.OrderUid
	ss.cache.AddToCache(*od)
	_, err = ss.db.SaveOrder(itemData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func (ss *StoreService) GetFromCacheByUID(id string) model.OrderData {
	return ss.cache.GetFromCache(id)
}

func (ss *StoreService) GetAllOrders() ([]model.DataItem, error) {
	di, err := ss.db.GetAllOrders()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return di, err
}

func (ss *StoreService) RestoreCache() error {
	dItems, err := ss.GetAllOrders()
	if dItems == nil {
		log.Println(err)
		return err
	}
	for _, dItem := range dItems {
		ss.cache.AddToCache(dItem.OrderData)
	}
	log.Println("--CACHE IS RESTORED--")
	return err
}
