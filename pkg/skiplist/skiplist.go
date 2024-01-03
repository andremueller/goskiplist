package skiplist

import (
	"cmp"
	"math/rand"
)

type SkipList[K cmp.Ordered, V any] struct {
	p        float64
	maxLevel int
	count    int
	head     *Node[K, V]
}

func NewSkipList[K cmp.Ordered, V any](p float64, maxLevel int) *SkipList[K, V] {
	var dummyKey K
	var dummyValue V
	return &SkipList[K, V]{
		p:        p,
		maxLevel: maxLevel,
		count:    0,
		head:     newNode[K, V](dummyKey, dummyValue, 1, maxLevel),
	}
}

func (s *SkipList[K, V]) Size() int {
	return s.count
}

func (s *SkipList[K, V]) Level() int {
	return s.head.Level()
}

func (s *SkipList[K, V]) randomLevel() int {
	level := 1

	for rand.Float64() < s.p && level < int(s.maxLevel) {
		level += 1
	}

	return level
}

func (s *SkipList[K, V]) Set(key K, value V) *Node[K, V] {
	update := make([]*Node[K, V], s.maxLevel)
	x := s.head
	for i := s.Level() - 1; i >= 0; i-- {
		for x.next[i] != nil && cmp.Less(x.next[i].key, key) {
			x = x.next[i]
		}
		update[i] = x
	}
	x = x.next[0]
	if x != nil && x.key == key {
		// override value
		x.Value = value
		return x
	}

	newLevel := s.randomLevel()

	if newLevel > s.Level() {
		for i := s.Level(); i < newLevel; i++ {
			update[i] = s.head
		}
		s.head.extendLevel(newLevel)
	}
	x = newNode[K, V](key, value, newLevel, newLevel)
	for i := 0; i < newLevel; i++ {
		x.next[i] = update[i].next[i]
		update[i].next[i] = x
		// TODO update dist on all levels
	}
	s.count++

	return x
}

func (s *SkipList[K, V]) Get(key K) *Node[K, V] {
	x := s.head
	for i := s.Level() - 1; i >= 0; i-- {
		for x.next[i] != nil && cmp.Less(x.next[i].key, key) {
			x = x.next[i]
		}
	}
	x = x.next[0]
	if x != nil && x.key == key {
		return x
	}
	return nil
}

// TODO dist not implement
func (s *SkipList[K, V]) GetByPos(k int) *Node[K, V] {
	if k >= s.count {
		return nil
	}
	x := s.head
	pos := 0
	for i := s.Level() - 1; i >= 0; i-- {
		for pos+x.dist[i] <= k {
			pos += x.dist[i]
			x = x.next[i]
		}
	}

	return x
}

// TODO InsertByPos
// TODO RemoveByPos

// TODO
// func (s *SkipList[K, V]) Remove(key K) bool {
// 	return false
// }
