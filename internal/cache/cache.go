package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	capacity int
	items    map[string]*list.Element
	ll       *list.List
	mu       sync.Mutex
}

type entry struct {
	key   string
	value interface{}
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		ll:       list.New(),
	}
}

