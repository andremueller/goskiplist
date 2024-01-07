package skiplist

import (
	"cmp"
	"fmt"
)

type Node[K cmp.Ordered, V any] struct {
	key   K
	Value V
	next  []*Node[K, V]
	dist  []int
}

func newNode[K cmp.Ordered, V any](key K, value V, level int, capacity int) *Node[K, V] {
	if capacity < level {
		capacity = level
	}
	return &Node[K, V]{
		key:   key,
		Value: value,
		next:  make([]*Node[K, V], level, capacity),
		dist:  make([]int, level, capacity),
	}
}

func (n *Node[K, V]) Key() K {
	return n.key
}

func (n *Node[K, V]) Level() int {
	return len(n.next)
}

func (n *Node[K, V]) Next() *Node[K, V] {
	if len(n.next) > 0 {
		return n.next[0]
	}
	return nil
}

func (n *Node[K, V]) extendLevel(newLevel int) {
	oldLevel := n.Level()
	if newLevel > oldLevel {
		n.next = n.next[:newLevel]
		n.dist = n.dist[:newLevel]
		for i := oldLevel; i < newLevel; i++ {
			n.next[i] = nil
			n.dist[i] = 0
		}
	}
}

func (n *Node[K, V]) shrinkLevel(newLevel int) {
	oldLevel := n.Level()
	if newLevel < oldLevel {
		n.next = n.next[:newLevel]
		n.dist = n.dist[:newLevel]
	}
}

func (n *Node[K, V]) String() string {
	return fmt.Sprintf("key: %v, value: %v | dist: %v", n.key, n.Value, n.dist)
}
