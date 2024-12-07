package cache

import (
	"caching-proxy/proxy/request"
	"sync"
)

type Request = request.Request

// thread-safe cache
type Cache struct {
	mtx    sync.RWMutex
	host   string
	backup string
	store  map[string]Request
}

// create new cache instance
func New(host, backup string) *Cache {
	// restore from backup
	reqs, _ := request.RestoreAll(backup)
	mp := make(map[string]Request)
	if r, ok := reqs[host]; ok {
		for _, v := range r {
			mp[v.Uri] = v
		}
	}
	return &Cache{
		store:  mp,
		host:   host,
		backup: backup,
	}
}

// add new or update existing
func (c *Cache) Put(key string, val Request) {
	c.mtx.Lock() // lock rw
	defer c.mtx.Unlock()
	val.Uri = key
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

func (c *Cache) Backup() error {
	c.mtx.Lock() // lock rw
	defer c.mtx.Unlock()
	reqs := make([]Request, len(c.store))
	for _, v := range c.store {
		reqs = append(reqs, v)
	}
	return request.BackupOne(c.backup, c.host, reqs)
}
