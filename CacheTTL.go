package CacheTTL

import (
	"sync"
	"time"
)

type storeField struct {
	value    interface{}
	deleteAt *time.Time
}

type Cache struct {
	data sync.Map
}

func New(searchExpiredInterval time.Duration) *Cache {
	cache := &Cache{
		data: sync.Map{},
	}
	go cache.backgroundCacheCleaner(searchExpiredInterval)
	return cache
}

func (c *Cache) backgroundCacheCleaner(searchExpiredInterval time.Duration) {
	for {
		<-time.Tick(searchExpiredInterval)
		c.data.Range(
			func(key, v interface{}) bool {
				vv, ok := v.(*storeField)
				if !ok {
					return true
				}

				// Если TTL не задан, то пропускаем
				if vv.deleteAt == nil {
					return true
				}

				if time.Now().After(*vv.deleteAt) {
					c.data.Delete(key)
				}
				return true
			})
	}
}

func (c *Cache) Set(key string, val interface{}, ttl time.Duration) {
	deleteAt := time.Now().Add(ttl)
	c.data.Store(key, &storeField{val, &deleteAt})
}

func (c *Cache) Get(key string) (result interface{}, ok bool) {
	load, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	vv, ok := load.(*storeField)
	if !ok {
		return nil, false
	}

	return vv.value, true
}

func (c *Cache) Delete(key string) {
	c.data.Delete(key)
}
