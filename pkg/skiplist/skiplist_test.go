package skiplist

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	key   int
	level int
	pos   int
}

// example Data of Figure 1 of "A skip List Cookbok"
var example1 = []testData{
	{3, 1, 0},
	{6, 4, 1},
	{7, 1, 2},
	{9, 2, 3},
	{12, 1, 4},
	{17, 2, 5},
	{19, 1, 6},
	{21, 1, 7},
	{25, 3, 8},
	{26, 1, 9},
}

var example2 = []testData{
	{3, 1, 0},
	{6, 4, 1},
	{7, 1, 2},
	{9, 2, 3},
	{12, 1, 4},
	{17, 2, 5},
	{19, 1, 6},
	{21, 1, 7},
	{26, 1, 9},
	{25, 3, 8},
}

// returns one level after the other instead of "random" levels for test use
func createPlayBackLevelFunc(data []testData) LevelFunc {
	pos := -1
	return func(p float64, maxLevel int) int {
		pos++
		if pos >= len(data) {
			log.Fatal("out of range in playback LevelFunc")
		}
		return data[pos].level
	}
}

func createSkipList(data []testData) *SkipList[int, int] {
	s := NewSkipList[int, int](WithLevelFunc[int, int](createPlayBackLevelFunc(data)))

	for i, x := range example1 {
		s.Set(x.key, i)
	}

	return s
}

func TestGetByPosWithFixed(t *testing.T) {
	data := example1
	s := NewSkipList[int, int](WithLevelFunc[int, int](createPlayBackLevelFunc(data)))

	fmt.Print(s.String())
	for i, x := range data {
		fmt.Printf("================== %d ==================== (%d)\n", i, x.key)
		_, pos, _ := s.Set(x.key, x.pos)
		fmt.Printf(" pos = %d\n", pos)
		fmt.Print(s.String())
	}

	for _, x := range data {
		n := s.GetByPos(x.pos)
		assert.NotNil(t, n)
		assert.Equal(t, x.key, n.Key())
	}
}

func TestGetByPosWithFixed2(t *testing.T) {
	data := make([]testData, len(example2))
	copy(data, example2)
	Shuffle(data)
	s := NewSkipList[int, int](WithLevelFunc[int, int](createPlayBackLevelFunc(data)))

	fmt.Print(s.String())
	for i, x := range data {
		fmt.Printf("================== %d ==================== (%d)\n", i, x.key)
		_, pos, _ := s.Set(x.key, x.pos)
		fmt.Printf(" pos = %d\n", pos)
		fmt.Print(s.String())
	}

	// GetByPos
	for _, x := range data {
		n := s.GetByPos(x.pos)
		assert.NotNil(t, n)
		assert.Equal(t, x.key, n.Key())
	}

	// Get
	for _, x := range data {
		n, pos := s.Get(x.key)
		assert.NotNil(t, n)
		assert.Equal(t, x.key, n.Key())
		assert.Equal(t, x.pos, pos)
	}

}

func Shuffle[V any](a []V) {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}

func makeRandomData(count int) []int {
	keys := make([]int, count)
	for i := 0; i < count; i++ {
		keys[i] = i
	}
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	return keys
}

func randomTest(t *testing.T, s *SkipList[int, int], count int) {
	keys := makeRandomData(count)

	// Set
	for i, k := range keys {
		assert.Equal(t, i, s.Size())
		node, pos := s.Get(k)
		assert.Nil(t, node)
		assert.Equal(t, InvalidPos, pos)
		s.Set(k, i)
		assert.Equal(t, i+1, s.Size())
	}

	// Get
	for i, k := range keys {
		x, pos := s.Get(k)
		assert.NotNil(t, x)
		assert.Equal(t, k, x.Key())
		assert.Equal(t, i, x.Value)
		assert.Equal(t, k, pos) // key will exactly match its position
	}

	// Remove
	n := s.Size()
	for i, k := range keys {
		x, pos := s.Remove(k)
		assert.NotNil(t, x)
		assert.True(t, pos >= 0)
		assert.Equal(t, k, x.Key())
		x2, pos2 := s.Get(k)
		assert.Nil(t, x2)
		assert.True(t, pos2 == InvalidPos)
		n--
		assert.Equal(t, n, s.Size())

		// check if all remaining elements are found
		for j := i + 1; j < len(keys); j++ {
			node, ppos := s.Get(keys[j])
			assert.NotNil(t, node)
			assert.True(t, ppos >= 0 && ppos < s.Size())

			node2 := s.GetByPos(ppos)
			assert.NotNil(t, node2)
			assert.Equal(t, node, node2)
		}
	}
}

func TestNewSkipList(t *testing.T) {
	s := NewSkipList[int, int]()
	randomTest(t, s, 100)
}

func randomPosTest(t *testing.T, s *SkipList[int, int], count int) {

	keys := makeRandomData(count)

	for i, k := range keys {
		assert.Equal(t, i, s.Size())
		node, pos := s.Get(k)
		assert.Nil(t, node)
		assert.Equal(t, InvalidPos, pos)
		s.Set(k, i)
		assert.Equal(t, i+1, s.Size())
	}

	for i, k := range keys {
		x := s.GetByPos(k)
		assert.NotNil(t, x)
		assert.Equal(t, k, x.Key())
		assert.Equal(t, i, x.Value)
	}
}

func TestGetByPos(t *testing.T) {
	s := NewSkipList[int, int]()
	randomPosTest(t, s, 100)
}
