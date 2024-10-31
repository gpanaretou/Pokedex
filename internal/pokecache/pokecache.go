package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mu       *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	ticker := time.NewTicker(interval)
	cleanUp := make(chan bool)

	cache := Cache{entries: map[string]cacheEntry{}, interval: interval, mu: &sync.Mutex{}}

	go func() {
		for {
			select {
			case <-cleanUp:
				return
			case <-ticker.C:
				cache.reapLoop()
			}
		}
	}()
	return cache
}

func (c Cache) Add(key string, val []byte) {
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries[key]
	if !ok {
		return []byte{}, false
	}

	return entry.val, true
}

func (c Cache) reapLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheEntries := c.entries
	for key, value := range cacheEntries {
		if time.Since(value.createdAt) > c.interval {
			delete(cacheEntries, key)
		}
	}
}
