package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewRow(t *testing.T) {
	data := RowData{"name": "Alice", "age": 30}
	row := NewRow(data)

	assert.Equal(t, "Alice", row.Data["name"])
	assert.Equal(t, 30, row.Data["age"])
	assert.False(t, row.IsSelected())
	assert.NotZero(t, row.ID())
}

func TestNewRowShallowCopy(t *testing.T) {
	data := RowData{"name": "Alice"}
	row := NewRow(data)

	// Mutating original data should not affect the row
	data["name"] = "Bob"
	assert.Equal(t, "Alice", row.Data["name"])
}

func TestNewRowNilData(t *testing.T) {
	row := NewRow(nil)
	assert.NotNil(t, row.Data)
	assert.Len(t, row.Data, 0)
}

func TestRowUniqueIDs(t *testing.T) {
	row1 := NewRow(nil)
	row2 := NewRow(nil)
	assert.NotEqual(t, row1.ID(), row2.ID())
}

func TestRowWithStyle(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true)
	row := NewRow(nil).WithStyle(style)
	assert.Equal(t, style, row.Style)
}

func TestRowSelected(t *testing.T) {
	row := NewRow(nil)
	assert.False(t, row.IsSelected())

	selected := row.Selected(true)
	assert.True(t, selected.IsSelected())
	assert.False(t, row.IsSelected(), "original should be unchanged")
}

func TestRowImmutability(t *testing.T) {
	original := NewRow(RowData{"x": 1})
	modified := original.WithStyle(lipgloss.NewStyle().Bold(true))

	assert.Equal(t, lipgloss.NewStyle(), original.Style)
	assert.NotEqual(t, original.Style, modified.Style)
}
