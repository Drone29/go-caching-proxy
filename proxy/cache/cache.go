package cache

import (
	"caching-proxy/proxy/request"
	"sync"
)

type Request = request.Request

// thread-safe cache
type Cache struct {
	mtx   sync.RWMutex
	store map[string]Request
}

// create new cache instance
func New() *Cache {
	return &Cache{
		store: make(map[string]Request),
	}
}

// add new or update existing
func (c *Cache) Put(key string, val Request) {
	c.mtx.Lock() // lock rw
	defer c.mtx.Unlock()
	c.store[key] = val
}

// retrieve from cache
func (c *Cache) Get(key string) (Request, bool) {
	c.mtx.RLock() // lock r
	defer c.mtx.RUnlock()
	res, ok := c.store[key]
	return res, ok
}

// delete from cache
func (c *Cache) Delete(key string) {
	c.mtx.Lock() // lock rw
	defer c.mtx.Unlock()
	delete(c.store, key)
}

// clear cache
func (c *Cache) Clear() {
	c.mtx.Lock() // lock rw
	defer c.mtx.Unlock()
	for k := range c.store {
		delete(c.store, k)
	}
}
