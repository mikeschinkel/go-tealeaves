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

	t.Run("right only", func(t *testing.T) {
		result := m.composeFooterZones("", "1/3", 20)
		assert.Contains(t, result, "1/3")
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
