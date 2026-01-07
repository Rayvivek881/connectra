package utilities

import (
	"testing"
	"time"
)

func TestTTLCache_GetSet(t *testing.T) {
	cache := NewTTLCache(5*time.Minute, 100)
	
	// Test Set and Get
	cache.Set("key1", "value1")
	
	value, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1 in cache")
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}
	
	// Test non-existent key
	_, found = cache.Get("nonexistent")
	if found {
		t.Error("Expected not to find nonexistent key")
	}
}

func TestTTLCache_Expiration(t *testing.T) {
	cache := NewTTLCache(100*time.Millisecond, 100)
	
	// Set a value with short TTL
	cache.Set("key1", "value1")
	
	// Should be found immediately
	_, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1 immediately after setting")
	}
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Should not be found after expiration
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be expired")
	}
}

func TestTTLCache_Delete(t *testing.T) {
	cache := NewTTLCache(5*time.Minute, 100)
	
	cache.Set("key1", "value1")
	cache.Delete("key1")
	
	_, found := cache.Get("key1")
	if found {
		t.Error("Expected key1 to be deleted")
	}
}

func TestTTLCache_Clear(t *testing.T) {
	cache := NewTTLCache(5*time.Minute, 100)
	
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Clear()
	
	if cache.Size() != 0 {
		t.Errorf("Expected cache size 0 after clear, got %d", cache.Size())
	}
}

func TestTTLCache_InvalidateByPrefix(t *testing.T) {
	cache := NewTTLCache(5*time.Minute, 100)
	
	cache.Set("filter:company:key1", "value1")
	cache.Set("filter:company:key2", "value2")
	cache.Set("filter:contact:key1", "value3")
	
	cache.InvalidateByPrefix("filter:company:")
	
	// Company keys should be gone
	_, found := cache.Get("filter:company:key1")
	if found {
		t.Error("Expected filter:company:key1 to be invalidated")
	}
	
	_, found = cache.Get("filter:company:key2")
	if found {
		t.Error("Expected filter:company:key2 to be invalidated")
	}
	
	// Contact key should still exist
	_, found = cache.Get("filter:contact:key1")
	if !found {
		t.Error("Expected filter:contact:key1 to still exist")
	}
}

func TestTTLCache_MaxSize(t *testing.T) {
	cache := NewTTLCache(5*time.Minute, 3)
	
	// Fill cache to max size
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	
	if cache.Size() != 3 {
		t.Errorf("Expected cache size 3, got %d", cache.Size())
	}
	
	// Adding one more should trigger eviction
	cache.Set("key4", "value4")
	
	// Size should still be at max
	if cache.Size() > 3 {
		t.Errorf("Expected cache size <= 3, got %d", cache.Size())
	}
}

func TestTTLCache_SetWithTTL(t *testing.T) {
	cache := NewTTLCache(5*time.Minute, 100)
	
	// Set with custom short TTL
	cache.SetWithTTL("key1", "value1", 50*time.Millisecond)
	
	// Should be found immediately
	_, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1 immediately")
	}
	
	// Wait for expiration
	time.Sleep(100 * time.Millisecond)
	
	// Should be expired
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be expired with custom TTL")
	}
}
