package skiplist

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randomTest(t *testing.T, p float64, maxLevel int, count int) {
	s := NewSkipList[int, int](p, maxLevel)

	keys := make([]int, count)
	for i := 0; i < count; i++ {
		keys[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	for i, k := range keys {
		assert.Equal(t, i, s.Size())
		s.Set(k, i)
		assert.Equal(t, i+1, s.Size())
		x := s.Get(k)
		assert.NotNil(t, x)

		assert.Equal(t, k, x.Key())
		assert.Equal(t, i, x.Value)
	}
}

func TestNewSkipList(t *testing.T) {
	randomTest(t, 0.5, 10, 100)
}
