package cache

import (
	"container/list"
	"sync"
)

// LRU Cache
type Cache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	items    map[K]*list.Element
	evict    *list.List
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element),
		evict:    list.New(),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the item already exists and update and move item to the front of the eviction list
	if el, ok := c.items[key]; ok {
		el.Value.(*entry[K, V]).value = value
		c.evict.MoveToFront(el)
		return
	}

	if len(c.items) >= c.capacity {
		// Remove the least recently used item
		c.evictOldest()
	}

	// Add new entry
	ent := &entry[K, V]{key, value}
	el := c.evict.PushFront(ent)
	c.items[key] = el
}

func (c *Cache[K, V]) evictOldest() {
	el := c.evict.Back()
	if el != nil {
		ent := el.Value.(*entry[K, V])
		delete(c.items, ent.key)
		c.evict.Remove(el)
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.evict.MoveToFront(el)
		return el.Value.(*entry[K, V]).value, true
	}

	var zero V
	return zero, false
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.evict.Remove(el)
		delete(c.items, key)
	}
}

func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
