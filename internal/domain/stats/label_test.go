package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLabel(t *testing.T) {
	label1 := newLabel("foo")
	assert.Equal(t, "foo", label1.name)
	assert.False(t, label1.negative)
	label2 := newLabel("-foo")
	assert.Equal(t, "foo", label2.name)
	assert.True(t, label2.negative)
}

func TestLabelGroup(t *testing.T) {
	g := newLabelGroup([]string{"foo", "foo2", "-bar"})
	assert.False(t, g.isFlagged())
	g.label("foo")
	assert.False(t, g.isFlagged())
	g.label("foo2")
	assert.True(t, g.isFlagged())
	g.label("bar")
	assert.False(t, g.isFlagged())
	g.unlabel("bar")
	assert.True(t, g.isFlagged())
	g.unlabel("foo")
	assert.False(t, g.isFlagged())
	g.label("not in group")
	g.unlabel("not in group")
}
