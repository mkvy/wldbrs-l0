package cache

import (
	"github.com/mkvy/wldbrs-l0/server-subscriber/model"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCacheInit(t *testing.T) {
	cache := CacheInit()
	require.NotNil(t, cache)
}

func TestCacheService_AddToCache(t *testing.T) {
	cache := CacheInit()
	require.NotNil(t, cache)
	fakeCache := make(CacheStore)
	text, err := os.ReadFile("../testdata/payload.json")
	if err != nil {
		t.Fatal("error reading testdata file")
	}
	od := new(model.OrderData)
	err = od.Scan(text)
	fakeCache[od.OrderUid] = *od
	cache.AddToCache(*od)
	realC := cache.CacheStore
	require.Equal(t, realC, fakeCache)
}

func TestCacheService_GetFromCache(t *testing.T) {
	cache := CacheInit()
	text, err := os.ReadFile("../testdata/payload.json")
	if err != nil {
		t.Fatal("error reading testdata file")
	}
	od := new(model.OrderData)
	err = od.Scan(text)
	cache.CacheStore[od.OrderUid] = *od
	fakeCache := make(CacheStore)
	fakeCache[od.OrderUid] = *od
	realOr := cache.GetFromCache(od.OrderUid)
	fakeOr := fakeCache[od.OrderUid]
	require.Equal(t, realOr, fakeOr)
}
