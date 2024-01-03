package skiplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomTest(t *testing.T, p float64, maxLevel int, count int) {
	s := NewSkipList[int, int](p, maxLevel)

	assert.Equal(t, uint64(0), s.Size())

	s.Set(127, 1)
	assert.Equal(t, 1, s.Size())

	x := s.Get(127)
	assert.NotNil(t, x)

	assert.Equal(t, 1, x.Value)
}

func TestNewSkipList(t *testing.T) {
	randomTest(t, 0.5, 10, 100)
}
