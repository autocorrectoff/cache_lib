package main

import (
	"fmt"
	"github.com/brankomiric/cache_lib/internal/cache"
)

func main() {
	// Init generic cache with capacity 3
	c := cache.NewCache[string, int](3)

	// Set some values within capacity
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	// An item outside capacity
	c.Set("d", 4)

	// Check the oldest item is gone
	_, ok := c.Get("a")
	fmt.Println("a exists:", ok)

	// Check random existing key
	val, ok := c.Get("b")
	if ok {
		fmt.Println("b value is:", val)
	}
}
