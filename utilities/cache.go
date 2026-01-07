package utilities

import (
	"sync"
	"time"
)

// CacheEntry represents a cached value with expiration time
type CacheEntry struct {
	Value      interface{}
	ExpiresAt  time.Time
}

// IsExpired checks if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// TTLCache is a thread-safe in-memory cache with TTL support
type TTLCache struct {
	mu        sync.RWMutex
	entries   map[string]*CacheEntry
	defaultTTL time.Duration
	maxSize   int
}

// NewTTLCache creates a new TTL cache instance
func NewTTLCache(defaultTTL time.Duration, maxSize int) *TTLCache {
	cache := &TTLCache{
		entries:   make(map[string]*CacheEntry),
		defaultTTL: defaultTTL,
		maxSize:   maxSize,
	}
	
	// Start background cleanup goroutine
	go cache.cleanup()
	
	return cache
}

// Get retrieves a value from the cache
func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}
	
	if entry.IsExpired() {
		// Don't delete here, let cleanup handle it
		return nil, false
	}
	
	return entry.Value, true
}

// Set stores a value in the cache with default TTL
func (c *TTLCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores a value in the cache with custom TTL
func (c *TTLCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if we need to evict entries (simple LRU-like eviction)
	if len(c.entries) >= c.maxSize {
		// Remove oldest expired entries first
		c.removeExpiredUnsafe()
		
		// If still at max size, remove a random entry (simple eviction)
		if len(c.entries) >= c.maxSize {
			// Remove first entry found (simple FIFO-like eviction)
			for k := range c.entries {
				delete(c.entries, k)
				break
			}
		}
	}
	
	c.entries[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a key from the cache
func (c *TTLCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Clear removes all entries from the cache
func (c *TTLCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

// InvalidateByPrefix removes all entries with keys starting with the prefix
func (c *TTLCache) InvalidateByPrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for key := range c.entries {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.entries, key)
		}
	}
}

// removeExpiredUnsafe removes expired entries (must be called with lock held)
func (c *TTLCache) removeExpiredUnsafe() {
	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.ExpiresAt) {
			delete(c.entries, key)
		}
	}
}

// cleanup periodically removes expired entries
func (c *TTLCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		c.removeExpiredUnsafe()
		c.mu.Unlock()
	}
}

// Size returns the current number of entries in the cache
func (c *TTLCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// Global cache instances for different data types
var (
	// FilterMetadataCache caches filter metadata (filters list)
	FilterMetadataCache = NewTTLCache(10*time.Minute, 100)
	
	// FilterDataCache caches filter data values
	FilterDataCache = NewTTLCache(5*time.Minute, 1000)
)
