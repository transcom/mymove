package migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {

	stack := NewStack()

	assert.True(t, stack.Empty())
	assert.Equal(t, 0, stack.Len())

	stack.Push("foo")

	assert.False(t, stack.Empty())
	assert.Equal(t, 1, stack.Len())
	assert.Equal(t, "foo", stack.Last())

	stack.Push("bar")

	assert.False(t, stack.Empty())
	assert.Equal(t, 2, stack.Len())
	assert.Equal(t, "bar", stack.Last())

	stack.Pop()

	assert.False(t, stack.Empty())
	assert.Equal(t, 1, stack.Len())
	assert.Equal(t, "foo", stack.Last())

	stack.Pop()

	assert.True(t, stack.Empty())
	assert.Equal(t, 0, stack.Len())

}
