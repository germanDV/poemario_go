package cache

import (
	"sync"

	"github.com/germandv/poemario/internal/list"
)

type LRUCache[K comparable, V any] struct {
	data map[K]*list.Node[K, V]
	ll   *list.LinkedList[K, V]
	max  int
	mu   sync.Mutex
}

func New[K comparable, V any](max int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		data: make(map[K]*list.Node[K, V]),
		ll:   &list.LinkedList[K, V]{},
		max:  max,
	}
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.data[key]
	if !exists {
		var zero V
		return zero, false
	}

	c.ll.MoveToHead(node)
	return node.Val, true
}

func (c *LRUCache[K, V]) Set(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.data[key]
	if exists {
		node.Val = val
		c.ll.MoveToHead(node)
	} else {
		node = c.ll.AddToHead(key, val)
		c.data[key] = node

		if len(c.data) > c.max {
			c.evict()
		}
	}
}

func (m *LRUCache[K, V]) evict() {
	removed := m.ll.RemoveTail()
	if removed != nil {
		delete(m.data, removed.Key)
	}
}
