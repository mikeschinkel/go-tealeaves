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

	assert.Len(t, m.GetVisibleRows(), 2)
}

func TestWithRowsResetsCache(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows([]Row{NewRow(RowData{"x": 1})})

	_ = m.GetVisibleRows() // populate cache

	m = m.WithRows([]Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	})

	assert.Len(t, m.GetVisibleRows(), 2)
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

func TestFocused(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).Focused(true)
	assert.True(t, m.GetFocused())
}

func TestWithCellCursorMode(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithCellCursorMode(true)
	assert.True(t, m.GetCellCursorMode())
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

	assert.Equal(t, 2, m.GetHighlightedRowIndex())
}

func TestWithHighlightedRowClamped(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		WithHighlightedRow(100)

	assert.Equal(t, 1, m.GetHighlightedRowIndex())
}

func TestWithHeaderVisibility(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithHeaderVisibility(false)
	assert.False(t, m.GetHeaderVisibility())
}

func TestWithFooterVisibility(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithFooterVisibility(false)
	assert.False(t, m.GetFooterVisibility())
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
	modified := original.Focused(true)

	assert.False(t, original.focused)
	assert.True(t, modified.focused)
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
