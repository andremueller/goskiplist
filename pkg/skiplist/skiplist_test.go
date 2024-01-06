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
}

// example Data of Figure 1 of "A skip List Cookbok"
var example1 = []testData{
	{3, 1},
	{6, 4},
	{7, 1},
	{9, 2},
	{12, 1},
	{17, 2},
	{19, 1},
	{21, 1},
	{25, 3},
	{26, 1},
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

func TestInsert(t *testing.T) {
	s := createSkipList(example1)
	fmt.Print(s.String())
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

	for i, k := range keys {
		assert.Equal(t, i, s.Size())
		node, pos := s.Get(k)
		assert.Nil(t, node)
		assert.Equal(t, InvalidPos, pos)
		s.Set(k, i)
		assert.Equal(t, i+1, s.Size())
	}

	for i, k := range keys {
		x, pos := s.Get(k)
		assert.NotNil(t, x)
		assert.Equal(t, k, x.Key())
		assert.Equal(t, i, x.Value)
		assert.Equal(t, k, pos) // key will exactly match its position
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
		x := s.GetByPos(i)
		assert.NotNil(t, x)
		assert.Equal(t, k, x.Key())
		assert.Equal(t, i, x.Value)
	}
}

func TestGetByPos(t *testing.T) {
	s := NewSkipList[int, int]()
	randomPosTest(t, s, 100)
}
