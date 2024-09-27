package cache

import (
	"testing"
)

func TestCacheSetAndGet(t *testing.T) {
	c := NewCache[string, int](2)

	// Set values
	c.Set("a", 1)
	c.Set("b", 2)

	// Test Get
	val, ok := c.Get("a")
	if !ok || val != 1 {
		t.Errorf("Expected 'a' to be 1, got %d", val)
	}

	val, ok = c.Get("b")
	if !ok || val != 2 {
		t.Errorf("Expected 'b' to be 2, got %d", val)
	}
}

func TestCacheEviction(t *testing.T) {
	c := NewCache[string, int](2)

	// Set values and trigger eviction
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3) // Evicts "a" because it's the least recently used

	// "a" should be evicted
	_, ok := c.Get("a")
	if ok {
		t.Errorf("Expected 'a' to be evicted, but it was found")
	}

	// "b" and "c" should still be present
	if val, ok := c.Get("b"); !ok || val != 2 {
		t.Errorf("Expected 'b' to be 2, got %d", val)
	}

	if val, ok := c.Get("c"); !ok || val != 3 {
		t.Errorf("Expected 'c' to be 3, got %d", val)
	}
}

func TestCacheUpdates(t *testing.T) {
	c := NewCache[string, int](2)

	// Set values
	c.Set("a", 1)
	c.Set("b", 2)

	// Update "a"
	c.Set("a", 10)

	val, ok := c.Get("a")
	if !ok || val != 10 {
		t.Errorf("Expected updated 'a' to be 10, got %d", val)
	}
}

func TestCacheDelete(t *testing.T) {
	c := NewCache[string, int](2)

	// Set and delete values
	c.Set("a", 1)
	c.Delete("a")

	// "a" should be deleted
	_, ok := c.Get("a")
	if ok {
		t.Errorf("Expected 'a' to be deleted, but it was found")
	}
}

func TestCacheEvictOldest(t *testing.T) {
	c := NewCache[string, int](2)

	// Set values
	c.Set("a", 1)
	c.Set("b", 2)

	// Access "a" to make it the most recently used
	c.Get("a")

	// Insert a new item, which should evict "b" because it's the least recently used
	c.Set("c", 3)

	// "b" should be evicted
	_, ok := c.Get("b")
	if ok {
		t.Errorf("Expected 'b' to be evicted, but it was found")
	}

	// "a" and "c" should still be present
	if val, ok := c.Get("a"); !ok || val != 1 {
		t.Errorf("Expected 'a' to be 1, got %d", val)
	}

	if val, ok := c.Get("c"); !ok || val != 3 {
		t.Errorf("Expected 'c' to be 3, got %d", val)
	}
}

func TestCacheConcurrency(t *testing.T) {
	c := NewCache[string, int](3)

	// Use goroutines to test concurrent access
	done := make(chan bool)

	// Writer
	go func() {
		for i := 0; i < 1000; i++ {
			c.Set(string(rune('a'+(i%26))), i)
		}
		done <- true
	}()

	// Reader
	go func() {
		for i := 0; i < 1000; i++ {
			c.Get(string(rune('a' + (i % 26))))
		}
		done <- true
	}()

	<-done
	<-done
}
