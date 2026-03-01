package teagrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeNaturalWidth(t *testing.T) {
	t.Run("fixed columns with rounded border", func(t *testing.T) {
		cols := []Column{
			NewColumn("a", "A", 10),
			NewColumn("b", "B", 20),
		}
		m := New(cols)
		// RenderWidth: (1+10+0) + (1+20+0) = 11 + 21 = 32
		// Outer border: 2
		// Inner divider: 1
		// Total: 32 + 2 + 1 = 35
		assert.Equal(t, 35, m.computeNaturalWidth())
	})

	t.Run("borderless", func(t *testing.T) {
		cols := []Column{
			NewColumn("a", "A", 10),
			NewColumn("b", "B", 20),
		}
		m := New(cols).WithBorder(Borderless())
		// RenderWidth: 11 + 21 = 32
		// No borders
		assert.Equal(t, 32, m.computeNaturalWidth())
	})

	t.Run("flex columns use minimum width", func(t *testing.T) {
		cols := []Column{
			NewColumn("a", "A", 10),
			NewFlexColumn("b", "B", 1),
		}
		m := New(cols)
		// Fixed: 11
		// Flex minimum: 1 + 1 + 0 = 2 (paddingLeft=1, content=1, paddingRight=0)
		// Outer: 2, Inner: 1
		// Total: 11 + 2 + 2 + 1 = 16
		assert.Equal(t, 16, m.computeNaturalWidth())
	})
}

func TestComputeTotalWidth(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 20),
	}
	m := New(cols)

	assert.Equal(t, 35, m.computeTotalWidth())
}

func TestUpdateColumnWidthsFlex(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewFlexColumn("b", "B", 1),
	}
	border := BorderRounded()

	updateColumnWidths(cols, 50, border)

	assert.Equal(t, 10, cols[0].Width(), "fixed column should keep its width")
	assert.Greater(t, cols[1].Width(), 0, "flex column should have resolved width")
}

func TestUpdateColumnWidthsMultipleFlex(t *testing.T) {
	cols := []Column{
		NewFlexColumn("a", "A", 1),
		NewFlexColumn("b", "B", 3),
	}
	border := BorderRounded()

	updateColumnWidths(cols, 44, border)

	// Total available: 44 - 2 (outer) - 1 (inner) = 41
	// Subtract padding: 41 - 1 - 0 - 1 - 0 = 39
	// Flex 1:3, so A gets ~10, B gets ~29
	assert.Greater(t, cols[0].Width(), 0)
	assert.Greater(t, cols[1].Width(), cols[0].Width())
}

func TestUpdateColumnWidthsNoFlex(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 20),
	}
	border := BorderRounded()

	updateColumnWidths(cols, 50, border)

	// No flex columns, widths should not change
	assert.Equal(t, 10, cols[0].Width())
	assert.Equal(t, 20, cols[1].Width())
}

func TestChromeHeight(t *testing.T) {
	t.Run("full chrome with rounded borders", func(t *testing.T) {
		m := New([]Column{NewColumn("x", "X", 10)})
		// outer=2, header=1, header_sep=1, footer=1, footer_sep=1 = 6
		assert.Equal(t, 6, m.chromeHeight())
	})

	t.Run("borderless", func(t *testing.T) {
		m := New([]Column{NewColumn("x", "X", 10)}).WithBorder(Borderless())
		// header=1, footer=1 = 2
		assert.Equal(t, 2, m.chromeHeight())
	})

	t.Run("no header", func(t *testing.T) {
		m := New([]Column{NewColumn("x", "X", 10)}).WithHeaderVisibility(false)
		// outer=2, footer=1, footer_sep=1 = 4
		assert.Equal(t, 4, m.chromeHeight())
	})

	t.Run("no footer", func(t *testing.T) {
		m := New([]Column{NewColumn("x", "X", 10)}).WithFooterVisibility(false)
		// outer=2, header=1, header_sep=1 = 4
		assert.Equal(t, 4, m.chromeHeight())
	})
}
