package teagrid

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewBasicRender(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10),
		NewColumn("age", "Age", 5),
	}
	rows := []Row{
		NewRow(RowData{"name": "Alice", "age": 30}),
		NewRow(RowData{"name": "Bob", "age": 25}),
	}

	m := NewGridModel(cols).WithRows(rows)
	output := m.render()

	assert.Contains(t, output, "Alice")
	assert.Contains(t, output, "Bob")
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "Age")
}

func TestViewBorderless(t *testing.T) {
	cols := []Column{
		NewColumn("x", "X", 5),
	}
	rows := []Row{
		NewRow(RowData{"x": "hello"}),
	}

	m := NewGridModel(cols).WithRows(rows).WithBorder(Borderless())
	output := m.render()

	// Should not contain any border characters
	assert.NotContains(t, output, "╭")
	assert.NotContains(t, output, "│")
	assert.NotContains(t, output, "╯")
	assert.Contains(t, output, "hello")
}

func TestViewNoRows(t *testing.T) {
	cols := []Column{
		NewColumn("x", "X", 5),
	}

	m := NewGridModel(cols)
	output := m.render()

	// Should at least render header
	assert.Contains(t, output, "X")
}

func TestViewFormatStringOnHeaders(t *testing.T) {
	cols := []Column{
		NewColumn("pct", "Win%", 8).WithFormatString("%.1f%%"),
	}
	rows := []Row{
		NewRow(RowData{"pct": 75.5}),
	}

	m := NewGridModel(cols).WithRows(rows)
	output := m.render()

	// Format string should apply to data
	assert.Contains(t, output, "75.5%")
}

func TestViewPadding(t *testing.T) {
	cols := []Column{
		NewColumn("x", "X", 5).WithPadding(2, 1),
	}
	rows := []Row{
		NewRow(RowData{"x": "hi"}),
	}

	m := NewGridModel(cols).WithRows(rows).WithBorder(Borderless()).WithHeaderVisibility(false)
	output := m.render()

	// With padding(2,1), content "hi" should have spaces around it
	// The cell should be: "  hi   " (2 left padding + "hi" padded to 5 + 1 right padding)
	assert.Contains(t, output, "  hi")
}

func TestViewAlignment(t *testing.T) {
	cols := []Column{
		NewColumn("x", "X", 10).WithAlignment(lipgloss.Right),
	}
	rows := []Row{
		NewRow(RowData{"x": "hi"}),
	}

	m := NewGridModel(cols).WithRows(rows).WithBorder(Borderless()).WithHeaderVisibility(false)
	output := m.render()

	// Right-aligned "hi" in 10-char width should have leading spaces
	assert.Contains(t, output, "        hi")
}

func TestViewHiddenHeader(t *testing.T) {
	cols := []Column{
		NewColumn("x", "X", 5),
	}
	rows := []Row{
		NewRow(RowData{"x": "data"}),
	}

	m := NewGridModel(cols).WithRows(rows).WithHeaderVisibility(false)
	output := m.render()

	// Header text should not appear
	lines := strings.Split(output, "\n")
	headerVisible := false
	for _, line := range lines {
		if strings.Contains(line, " X ") {
			headerVisible = true
		}
	}
	// In borderless mode, "X" might still appear as column title
	// With header hidden, the title row should not render
	require.False(t, headerVisible || strings.Contains(lines[0], " X "))
}

func TestViewMissingData(t *testing.T) {
	cols := []Column{
		NewColumn("x", "X", 10),
		NewColumn("y", "Y", 10),
	}
	rows := []Row{
		NewRow(RowData{"x": "present"}), // y is missing
	}

	m := NewGridModel(cols).WithRows(rows).WithMissingDataIndicator("N/A")
	output := m.render()

	assert.Contains(t, output, "present")
	assert.Contains(t, output, "N/A")
}

func TestViewCellValue(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true)
	cols := []Column{
		NewColumn("x", "X", 10),
	}
	rows := []Row{
		NewRow(RowData{
			"x": NewCellValue("styled", style),
		}),
	}

	m := NewGridModel(cols).WithRows(rows)
	output := m.render()

	assert.Contains(t, output, "styled")
}

func TestViewSpans(t *testing.T) {
	cols := []Column{
		NewColumn("x", "X", 20),
	}
	spans := []Span{
		NewSpan("hello", lipgloss.NewStyle()),
		NewSpan(" world", lipgloss.NewStyle()),
	}
	rows := []Row{
		NewRow(RowData{
			"x": NewCellValueWithSpans(spans, lipgloss.NewStyle()),
		}),
	}

	m := NewGridModel(cols).WithRows(rows)
	output := m.render()

	assert.Contains(t, output, "hello")
	assert.Contains(t, output, "world")
}

func TestPadOrTruncate(t *testing.T) {
	t.Run("left align pad", func(t *testing.T) {
		assert.Equal(t, "hi        ", padOrTruncate("hi", 10, lipgloss.Left))
	})

	t.Run("right align pad", func(t *testing.T) {
		assert.Equal(t, "        hi", padOrTruncate("hi", 10, lipgloss.Right))
	})

	t.Run("center align pad", func(t *testing.T) {
		result := padOrTruncate("hi", 10, lipgloss.Center)
		assert.Equal(t, 10, len(result))
		assert.Contains(t, result, "hi")
	})

	t.Run("exact width", func(t *testing.T) {
		assert.Equal(t, "hello", padOrTruncate("hello", 5, lipgloss.Left))
	})
}
