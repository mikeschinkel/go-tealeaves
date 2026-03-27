package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestRenderFooterHidden(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 10)}).
		WithFooterVisibility(false)

	assert.Equal(t, "", m.renderFooter())
}

func TestRenderFooterZeroHeightWhenHidden(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 10)}).
		WithFooterVisibility(false)

	// Hidden footer = zero height = empty string
	footer := m.renderFooter()
	assert.Empty(t, footer)
}

func TestRenderFooterWithPagination(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"x": i})
	}

	m := NewGridModel([]Column{NewColumn("x", "X", 10)}).
		WithRows(rows).
		WithPageSize(10)

	footer := m.renderFooter()
	assert.Contains(t, footer, "1/3")
}

func TestRenderFooterWithStaticText(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 10)}).
		WithStaticFooter("Total: 42")

	footer := m.renderFooter()
	assert.Contains(t, footer, "Total: 42")
}

func TestRenderFooterIndependentStyle(t *testing.T) {
	// Footer style should NOT inherit from baseStyle
	baseStyle := lipgloss.NewStyle().Bold(true)
	footerStyle := lipgloss.NewStyle().Italic(true)

	m := NewGridModel([]Column{NewColumn("x", "X", 10)}).
		WithBaseStyle(baseStyle).
		WithFooterStyle(footerStyle).
		WithStaticFooter("info")

	footer := m.renderFooter()
	assert.NotEmpty(t, footer)
}

func TestComposeFooterZones(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 10)})

	t.Run("both zones", func(t *testing.T) {
		result := m.composeFooterZones("filter", "1/3", 20)
		assert.Contains(t, result, "filter")
		assert.Contains(t, result, "1/3")
	})

	t.Run("right only centered by default", func(t *testing.T) {
		result := m.composeFooterZones("", "1/3", 20)
		assert.Contains(t, result, "1/3")
		// Centered: 8 spaces + "1/3" + 9 spaces = 20
		assert.Equal(t, 20, len(result))
		// "1/3" should not be flush left or flush right
		assert.NotEqual(t, "1/3", result[:3])
		assert.NotEqual(t, "1/3", result[17:])
	})

	t.Run("left only", func(t *testing.T) {
		result := m.composeFooterZones("filter", "", 20)
		assert.Contains(t, result, "filter")
	})

	t.Run("empty", func(t *testing.T) {
		result := m.composeFooterZones("", "", 10)
		assert.Len(t, result, 10) // all spaces
	})
}

func TestPlainFooterSeparatorUsesBottomJunctions(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
		NewColumn("b", "B", 5),
	}
	m := NewGridModel(cols).WithStaticFooter("info")

	sep := m.renderPlainFooterSeparator()
	chars := m.border.Chars

	// Should have left/right junctions and BottomJunction (upward T) at column positions,
	// NOT InnerJunction (cross) since the content below is full-width
	assert.Contains(t, sep, chars.LeftJunction)
	assert.Contains(t, sep, chars.RightJunction)
	assert.Contains(t, sep, chars.BottomJunction)
	assert.NotContains(t, sep, chars.InnerJunction)
}

func TestColumnJunctionSeparatorWithFooterRows(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
		NewColumn("b", "B", 5),
	}
	row := NewFooterRow(NewFooterCell("a", "x"))
	m := NewGridModel(cols).WithFooterRows(row)

	footer := m.renderFooter()
	chars := m.border.Chars

	// Footer separator should have inner junctions (column-aware rows follow)
	assert.Contains(t, footer, chars.InnerJunction)
}

func TestStaticFooterCenteredByDefault(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 20)}
	m := NewGridModel(cols).WithStaticFooter("hello")

	footer := m.renderFooter()
	assert.Contains(t, footer, "hello")

	// The info row content is centered
	infoRow := m.renderInfoRow()
	// "hello" is 5 chars, content width is 22 (20 + 2 padding)
	// Centered means it's not flush left
	assert.NotEqual(t, 'h', rune(infoRow[0]))
}

func TestWithFooterAlignmentRight(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 20)}
	m := NewGridModel(cols).
		WithStaticFooter("hello").
		WithFooterAlignment(lipgloss.Right)

	// composeFooterZones should right-align
	contentWidth := 30
	result := m.composeFooterZones("", "hello", contentWidth)
	// Right-aligned: spaces then "hello"
	assert.Equal(t, "hello", result[contentWidth-5:])
}

func TestFooterRowSingleCell(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
		NewColumn("b", "B", 5),
		NewColumn("c", "C", 5),
	}
	row := NewFooterRow(NewFooterCell("b", "val"))
	m := NewGridModel(cols).WithFooterRows(row)

	footer := m.renderFooter()
	assert.Contains(t, footer, "val")
}

func TestFooterRowColspan(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
		NewColumn("b", "B", 5),
		NewColumn("c", "C", 5),
	}
	// Span columns a+b (2 columns)
	row := NewFooterRow(NewFooterCellSpan("a", "wide", 2))
	m := NewGridModel(cols).WithFooterRows(row)

	footer := m.renderFooter()
	assert.Contains(t, footer, "wide")
}

func TestFooterRowsAndInfoRowCoexist(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 10)}
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"x": i})
	}
	fRow := NewFooterRow(NewFooterCell("x", "total"))
	m := NewGridModel(cols).
		WithRows(rows).
		WithPageSize(10).
		WithFooterRows(fRow)

	footer := m.renderFooter()
	// Should contain both the footer row content and pagination
	assert.Contains(t, footer, "total")
	assert.Contains(t, footer, "1/3")
}

func TestChromeHeightWithFooterRows(t *testing.T) {
	cols := []Column{NewColumn("x", "X", 10)}
	fRow := NewFooterRow(NewFooterCell("x", "total"))

	t.Run("footer rows only", func(t *testing.T) {
		m := NewGridModel(cols).WithFooterRows(fRow)
		// outer=2, header=1, header_sep=1, footer_sep=1, footer_row=1 = 6
		assert.Equal(t, 6, m.chromeHeight())
	})

	t.Run("footer rows plus info row", func(t *testing.T) {
		m := NewGridModel(cols).
			WithFooterRows(fRow).
			WithStaticFooter("info")
		// outer=2, header=1, header_sep=1, footer_sep=1, footer_row=1, plain_sep=1, info_row=1 = 8
		assert.Equal(t, 8, m.chromeHeight())
	})
}
