package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu       sync.Mutex
	entry    map[string]CacheEntry
	interval time.Duration
}

type CacheEntry struct {
	createdAt time.Time
	value     []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entry:    make(map[string]CacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) AddEntry(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = CacheEntry{
		createdAt: time.Now(),
		value:     value,
	}
}

func (c *Cache) DeleteEntry(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entry, key)
}

func (c *Cache) GetEntry(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entry[key]
	if !ok {
		return nil, false
	}
	return entry.value, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entry {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.entry, key)
		}
	}
}
