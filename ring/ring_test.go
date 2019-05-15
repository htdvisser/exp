package ring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRing(t *testing.T) {
	assert := assert.New(t)

	r := New(4)

	last, ok := r.Last()
	assert.False(ok)
	assert.Equal(zero, last)

	r.Add(1)
	r.Add(2)

	last, ok = r.Last()
	assert.Equal(Entry(2), last)
	assert.Equal([]Entry{1, 2}, r.All())

	r.Add(3)
	r.Add(4)

	last, ok = r.Last()
	assert.Equal(Entry(4), last)
	assert.Equal([]Entry{1, 2, 3, 4}, r.All())

	r.Add(5)
	r.Add(6)

	last, ok = r.Last()
	assert.Equal(Entry(6), last)
	assert.Equal([]Entry{3, 4, 5, 6}, r.All())
}
