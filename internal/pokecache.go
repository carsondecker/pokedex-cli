package internal

import (
	"sync"
	"time"
)

type Cache struct {
	caches map[string]cacheEntry
	mu     sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		caches: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.caches[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.caches[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		c.mu.Lock()
		for key, entry := range c.caches {
			if time.Since(entry.createdAt) > interval {
				delete(c.caches, key)
			}
		}
		c.mu.Unlock()
	}
}
