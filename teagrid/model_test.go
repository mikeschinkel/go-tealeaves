package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10),
		NewColumn("age", "Age", 5),
	}
	m := New(cols)

	assert.Len(t, m.columns, 2)
	assert.Equal(t, "name", m.columns[0].Key())
	assert.True(t, m.headerVisible)
	assert.True(t, m.footerVisible)
	assert.False(t, m.focused)
	assert.False(t, m.cellCursorMode)
}

func TestNewDefaultAlignment(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 10)}
	m := New(cols)

	// v0.2.0: baseStyle should NOT have right-align (fixes v0.1.0 #1)
	assert.Equal(t, lipgloss.NewStyle(), m.baseStyle)
}

func TestNewBorderDefault(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 10)}
	m := New(cols)

	// Default is rounded borders
	assert.True(t, m.border.HasOuterBorder())
	assert.Equal(t, "╭", m.border.Chars.TopLeft)
}

func TestNewImmutableColumns(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 10)}
	m := New(cols)

	// Mutating original slice should not affect model
	cols[0] = NewColumn("y", "Y", 20)
	assert.Equal(t, "x", m.columns[0].Key())
}

func TestSetSize(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewFlexColumn("b", "B", 1),
	}
	m := New(cols).SetSize(80, 24)

	assert.Equal(t, 80, m.viewportWidth)
	assert.Equal(t, 24, m.viewportHeight)

	// Flex column should have been resolved
	assert.Greater(t, m.columns[1].Width(), 0)
}

func TestSetSizeAutoPageSize(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 10)}
	m := New(cols).SetSize(80, 24)

	// With rounded borders (outer=2, header=1, header_sep=1, footer=1, footer_sep=1)
	// Chrome = 6, so pageSize = 24 - 6 = 18
	assert.Greater(t, m.pageSize, 0)
	assert.LessOrEqual(t, m.pageSize, 24)
}

func TestSetSizeBorderless(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 10)}
	m := New(cols).WithBorder(Borderless()).SetSize(80, 24)

	// Borderless: chrome = 0 (no outer, no header sep, no footer sep)
	// But header and footer rows still count
	assert.Greater(t, m.pageSize, 0)
}

func TestInit(t *testing.T) {
	m := New([]Column{NewColumn("x", "X", 10)})
	assert.Nil(t, m.Init())
}
