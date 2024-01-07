package skiplist

import (
	"cmp"
	"fmt"
	"log"
	"math/rand"
)

type LevelFunc func(p float64, maxLevel int) int

const MaxLevel = 512
const DefaultMaxLevel = 64
const DefaultProbability = 0.5

func defaultLevelFunc(p float64, maxLevel int) int {
	level := 1

	for rand.Float64() < p && level < maxLevel {
		level++
	}

	return level
}

type SkipList[K cmp.Ordered, V any] struct {
	p         float64     // probability for increasing the level of the skip list
	maxLevel  int         // maximum levels of the skip list
	count     int         // count is the number of elements in the skip list
	head      *Node[K, V] // the head node of the skip list
	levelFunc LevelFunc
}

type skipListOption[K cmp.Ordered, V any] func(*SkipList[K, V])

func WithLevelFunc[K cmp.Ordered, V any](levelFunc LevelFunc) skipListOption[K, V] {
	return func(s *SkipList[K, V]) {
		s.levelFunc = levelFunc
	}
}

func WithMaxLevel[K cmp.Ordered, V any](maxLevel int) skipListOption[K, V] {
	if maxLevel < 1 || maxLevel > MaxLevel {
		log.Panic("Parameter maxLevel out of range (must be >=1 and <= MaxLevel)")
	}
	return func(s *SkipList[K, V]) {
		s.maxLevel = maxLevel
	}
}

func WithProbability[K cmp.Ordered, V any](prob float64) skipListOption[K, V] {
	if prob < 0.01 || prob > 0.99 {
		log.Panic("Parameter probability out of range (must be >= 0.01 and <= 0.99)")
	}
	return func(s *SkipList[K, V]) {
		s.p = prob
	}
}

func NewSkipList[K cmp.Ordered, V any](options ...skipListOption[K, V]) *SkipList[K, V] {
	var dummyKey K
	var dummyValue V
	s := &SkipList[K, V]{
		p:         DefaultProbability,
		maxLevel:  DefaultMaxLevel,
		count:     0,
		levelFunc: defaultLevelFunc,
	}

	for _, opt := range options {
		opt(s)
	}

	s.head = newNode[K, V](dummyKey, dummyValue, 0, s.maxLevel)
	return s
}

func (s *SkipList[K, V]) First() *Node[K, V] {
	return s.head.Next()
}

func (s *SkipList[K, V]) Size() int {
	return s.count
}

func (s *SkipList[K, V]) Level() int {
	return s.head.Level()
}

func (s *SkipList[K, V]) randomLevel() int {
	return s.levelFunc(s.p, s.maxLevel)

}

// Set sets the value `value` of a key `key` within the skip list.
// Replaces the value if the key was already added to the set or inserts the key if not.
// Return a reference to the node and its current position 0...n-1 within the skip list.
// The bool value is true, if a new node was created and false if the value was overridden.
func (s *SkipList[K, V]) Set(key K, value V) (*Node[K, V], int, bool) {
	update := make([]*Node[K, V], s.Level(), s.maxLevel)
	updatePos := make([]int, s.Level(), s.maxLevel)
	x := s.head
	pos := -1 // the head has position -1, the first element 0
	for i := s.Level() - 1; i >= 0; i-- {
		for x.next[i] != nil && cmp.Less(x.next[i].key, key) {
			pos += x.dist[i]
			x = x.next[i]
		}
		update[i] = x
		updatePos[i] = pos
	}
	if len(x.next) > 0 && x.next[0] != nil && x.next[0].key == key {
		// key already exists: override value
		x = x.next[0]
		x.Value = value
		return x, pos, false
	}

	// now x.key shall be smaller than key
	newLevel := s.randomLevel()

	if newLevel > s.Level() {
		update = update[:newLevel]
		updatePos = updatePos[:newLevel]
		oldLevel := s.Level()
		s.head.extendLevel(newLevel)
		for i := oldLevel; i < newLevel; i++ {
			update[i] = s.head
			updatePos[i] = -1
			s.head.dist[i] = s.Size() + 1
		}
	}
	x = newNode[K, V](key, value, newLevel, newLevel)
	for i := 0; i < s.Level(); i++ {
		if i >= newLevel {
			update[i].dist[i]++
		} else {
			x.next[i] = update[i].next[i]
			update[i].next[i] = x
			delta := pos - updatePos[i]
			x.dist[i] = update[i].dist[i] - delta
			update[i].dist[i] = delta + 1
		}
	}

	s.count++

	return x, pos + 1, true
}

const InvalidPos = -1

// Get returns the node matching the searched key or nil if it was not found. The second return argument is the
// position 0...n-1 of the key
func (s *SkipList[K, V]) Get(key K) (*Node[K, V], int) {
	x := s.head
	pos := -1
	for i := s.Level() - 1; i >= 0; i-- {
		for x.next[i] != nil && x.next[i].key < key {
			pos += x.dist[i]
			x = x.next[i]
		}
	}
	if len(x.next) > 0 {
		x = x.next[0]
		pos++
		if x != nil && x.key == key {
			return x, pos
		}
	}
	return nil, InvalidPos
}

func (s *SkipList[K, V]) GetByPos(k int) *Node[K, V] {
	if k < 0 || k >= s.count {
		return nil
	}
	x := s.head
	pos := -1
	for i := s.Level() - 1; i >= 0; i-- {
		for x.next[i] != nil && pos+x.dist[i] <= k {
			pos += x.dist[i]
			x = x.next[i]
		}
	}

	return x
}

func (s *SkipList[K, V]) Remove(key K) (*Node[K, V], int) {
	update := make([]*Node[K, V], s.Level(), s.maxLevel)
	updatePos := make([]int, s.Level(), s.maxLevel)
	x := s.head
	pos := -1 // the head has position -1, the first element 0
	for i := s.Level() - 1; i >= 0; i-- {
		for x.next[i] != nil && cmp.Less(x.next[i].key, key) {
			pos += x.dist[i]
			x = x.next[i]
		}
		update[i] = x
		updatePos[i] = pos
	}
	if len(x.next) > 0 && x.next[0] != nil && x.next[0].key == key {
		// key found
		x = x.next[0]
		pos++

		// remove node from list
		for i := 0; i < s.Level(); i++ {
			if update[i].next[i] == x {
				update[i].next[i] = x.next[i]
				update[i].dist[i] += x.dist[i] - 1
			} else {
				update[i].dist[i]--
			}
		}

		// adapt level
		newLevel := s.Level()
		for newLevel > 0 && s.head.next[newLevel-1] == nil {
			newLevel--
		}
		s.head.shrinkLevel(newLevel)
		s.count--

		return x, pos
	}
	return nil, InvalidPos
}

func (s *SkipList[K, V]) RemoveByPos(k int) *Node[K, V] {
	if k < 0 || k >= s.count {
		return nil
	}

	update := make([]*Node[K, V], s.Level(), s.maxLevel)
	updatePos := make([]int, s.Level(), s.maxLevel)
	x := s.head
	pos := -1 // the head has position -1, the first element 0

	for i := s.Level() - 1; i >= 0; i-- {
		for x.next[i] != nil && pos+x.dist[i] < k {
			pos += x.dist[i]
			x = x.next[i]
		}
		update[i] = x
		updatePos[i] = pos
	}
	// remove node from list
	pos++
	x = x.Next()
	for i := 0; i < s.Level(); i++ {
		if update[i].next[i] == x {
			update[i].next[i] = x.next[i]
			update[i].dist[i] += x.dist[i] - 1
		} else {
			update[i].dist[i]--
		}
	}

	// adapt level
	newLevel := s.Level()
	for newLevel > 0 && s.head.next[newLevel-1] == nil {
		newLevel--
	}
	s.head.shrinkLevel(newLevel)
	s.count--

	return x
}

func (s *SkipList[K, V]) String() string {
	str := fmt.Sprintf("n=%d L=%d\n", s.Size(), s.Level())

	x := s.head
	for x != nil {
		str += x.String() + "\n"
		if len(x.next) > 0 {
			x = x.next[0]
		} else {
			x = nil
		}
	}
	return str
}
