package CacheTTL

import (
	"fmt"
	"sync"
	"time"
)

type value struct {
	value interface{}
	ttl   *time.Time
}

type Cache struct {
	data sync.Map
}

func New() *Cache {
	cache := &Cache{
		data: sync.Map{},
	}

	go cache.backgroundCacheCleaner()
	return cache
}

// background goroutine to clean up expired keys in the Cache
func (c *Cache) backgroundCacheCleaner() {
	for {
		<-time.Tick(time.Second * 1)
		fmt.Println("Searching for expired keys")
		c.data.Range(
			func(key, v interface{}) bool {
				vv, ok := v.(*value)
				if !ok {
					return true
				}

				if vv.ttl == nil {
					return true
				}

				if time.Now().After(*vv.ttl) {
					c.data.Delete(key)
					fmt.Println("Deleting expired record with key:", key)
				}
				return true
			})
	}
}

func (c *Cache) Set(key string, val interface{}, ttl time.Duration) {
	t := time.Now().Add(ttl)
	c.data.Store(key, &value{val, &t})
}

func (c *Cache) Get(key string) (result interface{}, ok bool) {
	load, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	vv, ok := load.(*value)
	if !ok {
		return nil, false
	}

	return vv.value, true
}

func (c *Cache) Delete(key string) {
	c.data.Delete(key)
}
