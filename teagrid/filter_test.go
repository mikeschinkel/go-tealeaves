package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestFilterFuncContains(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
	}

	t.Run("matches substring", func(t *testing.T) {
		row := NewRow(RowData{"name": "Alice"})
		assert.True(t, filterFuncContains(FilterFuncInput{
			Columns: cols, Row: row, Filter: "lic",
		}))
	})

	t.Run("case insensitive", func(t *testing.T) {
		row := NewRow(RowData{"name": "Alice"})
		assert.True(t, filterFuncContains(FilterFuncInput{
			Columns: cols, Row: row, Filter: "ALICE",
		}))
	})

	t.Run("no match", func(t *testing.T) {
		row := NewRow(RowData{"name": "Alice"})
		assert.False(t, filterFuncContains(FilterFuncInput{
			Columns: cols, Row: row, Filter: "Bob",
		}))
	})

	t.Run("empty filter matches all", func(t *testing.T) {
		row := NewRow(RowData{"name": "Alice"})
		assert.True(t, filterFuncContains(FilterFuncInput{
			Columns: cols, Row: row, Filter: "",
		}))
	})

	t.Run("CellValue data extracted", func(t *testing.T) {
		row := NewRow(RowData{
			"name": NewCellValue("Alice", lipgloss.NewStyle()),
		})
		assert.True(t, filterFuncContains(FilterFuncInput{
			Columns: cols, Row: row, Filter: "ali",
		}))
	})
}

func TestFilterFuncFuzzy(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
	}

	t.Run("subsequence match", func(t *testing.T) {
		row := NewRow(RowData{"name": "Alice"})
		assert.True(t, filterFuncFuzzy(FilterFuncInput{
			Columns: cols, Row: row, Filter: "ace",
		}))
	})

	t.Run("no match", func(t *testing.T) {
		row := NewRow(RowData{"name": "Alice"})
		assert.False(t, filterFuncFuzzy(FilterFuncInput{
			Columns: cols, Row: row, Filter: "xyz",
		}))
	})
}

func TestFuzzySubsequenceMatch(t *testing.T) {
	assert.True(t, fuzzySubsequenceMatch("hello world", "hlo"))
	assert.True(t, fuzzySubsequenceMatch("hello world", ""))
	assert.False(t, fuzzySubsequenceMatch("hello", "xyz"))
	assert.True(t, fuzzySubsequenceMatch("hello", "hello"))
}

func TestGetFilteredRows(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
	}
	rows := []Row{
		NewRow(RowData{"name": "Alice"}),
		NewRow(RowData{"name": "Bob"}),
		NewRow(RowData{"name": "Charlie"}),
	}

	m := New(cols).WithRows(rows).Filtered(true)
	m.filterTextInput.SetValue("ali")

	filtered := m.getFilteredRows(rows)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "Alice", filtered[0].Data["name"])
}
