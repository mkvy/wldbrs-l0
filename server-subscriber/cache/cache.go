package cache

import (
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
	"log"
)

type CacheStore map[string]model.OrderData

type CacheService struct {
	CacheStore CacheStore
}

func CacheInit() *CacheService {
	Cs := make(CacheStore)
	CacheService := CacheService{
		CacheStore: Cs,
	}
	return &CacheService
}

func (Cservice *CacheService) AddToCache(data model.OrderData) {
	Cservice.CacheStore[data.OrderUid] = data
	log.Println("new data in cache stored: ", data)
}

func (Cservice *CacheService) GetFromCache(order_uid string) model.OrderData {
	return Cservice.CacheStore[order_uid]
}
