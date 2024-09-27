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

	// Insert the new item
	ent := &entry[K, V]{key, value}
	el := c.evict.PushFront(ent)
	c.items[key] = el
}


