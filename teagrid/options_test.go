package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestWithRows(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithRows(rows)

	assert.Len(t, m.VisibleRows(), 2)
}

func TestWithRowsResetsCache(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows([]Row{NewRow(RowData{"x": 1})})

	_ = m.VisibleRows() // populate cache

	m = m.WithRows([]Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	})

	assert.Len(t, m.VisibleRows(), 2)
}

func TestWithBaseStyle(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true)
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithBaseStyle(style)
	assert.Equal(t, style, m.baseStyle)
}

func TestWithBorder(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithBorder(Borderless())
	assert.False(t, m.border.HasOuterBorder())
}

func TestWithFocused(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithFocused(true)
	assert.True(t, m.IsFocused())
}

func TestWithColCursorMode(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithColCursorMode(true)
	assert.True(t, m.ColCursorMode())
}

func TestWithSelectableRows(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithSelectableRows(true)
	assert.True(t, m.selectableRows)
	// v0.2.0: no auto-added column
	assert.Len(t, m.columns, 1)
}

func TestWithHighlightedRow(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
		NewRow(RowData{"x": 3}),
	}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		WithHighlightedRow(2)

	assert.Equal(t, 2, m.HighlightedRowIndex())
}

func TestWithHighlightedRowClamped(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		WithHighlightedRow(100)

	assert.Equal(t, 1, m.HighlightedRowIndex())
}

func TestWithHeaderVisibility(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithHeaderVisibility(false)
	assert.False(t, m.IsHeaderVisible())
}

func TestWithFooterVisibility(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithFooterVisibility(false)
	assert.False(t, m.IsFooterVisible())
}

func TestWithMetadata(t *testing.T) {
	meta := map[string]any{"theme": "dark"}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithMetadata(meta)
	assert.Equal(t, "dark", m.metadata["theme"])
}

func TestWithOverflowIndicator(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithOverflowIndicator(true)
	assert.True(t, m.overflowIndicator)
}

func TestImmutability(t *testing.T) {
	original := NewGridModel([]Column{NewColumn("x", "X", 5)})
	modified := original.WithFocused(true)

	assert.False(t, original.focused)
	assert.True(t, modified.focused)
}

func TestWithCellPadding(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 20),
	}
	m := NewGridModel(cols).WithCellPadding(2, 3)

	for _, col := range m.columns {
		assert.Equal(t, 2, col.PaddingLeft())
		assert.Equal(t, 3, col.PaddingRight())
	}
}

func TestWithCellPaddingZero(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 20),
	}
	m := NewGridModel(cols).WithCellPadding(0, 0)

	for _, col := range m.columns {
		assert.Equal(t, 0, col.PaddingLeft())
		assert.Equal(t, 0, col.PaddingRight())
	}
}

func TestWithFillWidth(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)})
	assert.True(t, m.FillWidth(), "fill-width should be enabled by default")

	m = m.WithFillWidth(true)
	assert.True(t, m.FillWidth())

	m = m.WithFillWidth(false)
	assert.False(t, m.FillWidth())
}

func TestWithEditableStub(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithEditable(true)
	assert.True(t, m.editable)
}

func TestWithCellValidatorStub(t *testing.T) {
	validator := func(columnKey string, value any) error { return nil }
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithCellValidator(validator)
	assert.NotNil(t, m.cellValidator)
}
