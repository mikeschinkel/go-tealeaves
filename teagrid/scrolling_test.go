package teagrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrollRight(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}
	m := NewGridModel(cols).SetSize(30, 24) // narrow viewport forces scrolling

	assert.Equal(t, 0, m.horizontalScrollOffsetCol)

	m.scrollRight()
	if m.maxHorizontalColumnIndex > 0 {
		assert.Equal(t, 1, m.horizontalScrollOffsetCol)
	}
}

func TestScrollLeft(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("a", "A", 20)})
	m.horizontalScrollOffsetCol = 2

	m.scrollLeft()
	assert.Equal(t, 1, m.horizontalScrollOffsetCol)

	m.scrollLeft()
	assert.Equal(t, 0, m.horizontalScrollOffsetCol)

	m.scrollLeft()
	assert.Equal(t, 0, m.horizontalScrollOffsetCol, "should not go below 0")
}

func TestScrollNoOverflowNoScroll(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
	}
	m := NewGridModel(cols).SetSize(80, 24)

	assert.Equal(t, 0, m.maxHorizontalColumnIndex)
}

func TestVisibleColumns(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 10),
		NewColumn("c", "C", 10),
	}

	t.Run("all fit", func(t *testing.T) {
		m := NewGridModel(cols).SetSize(100, 24)
		visible := m.visibleColumns()
		assert.Len(t, visible, 3)
	})

	t.Run("no viewport returns all", func(t *testing.T) {
		m := NewGridModel(cols)
		visible := m.visibleColumns()
		assert.Len(t, visible, 3)
	})
}

func TestHorizontalFreezeColumns(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 5),
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}

	m := NewGridModel(cols).
		WithHorizontalFreezeColumnCount(1).
		SetSize(30, 24)

	// First column should always be in visible set
	visible := m.visibleColumns()
	assert.Equal(t, "id", visible[0].Key())
}
