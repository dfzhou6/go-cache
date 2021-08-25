package dfcache

import (
	"dfcache/lru"
	"sync"
)

type cache struct {
	mu       sync.Mutex
	lru      *lru.LRU
	maxBytes int64
}

// Add key pair
func (c *cache) add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.maxBytes, nil)
	}

	c.lru.Add(key, val)
}

// Get key pair
func (c *cache) get(key string) (val ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
