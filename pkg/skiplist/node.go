package skiplist

import "cmp"

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

func (n *Node[K, V]) extendLevel(newLevel int) {
	if newLevel > n.Level() {
		n.next = n.next[:newLevel]
		n.dist = n.dist[:newLevel]
	}
}
