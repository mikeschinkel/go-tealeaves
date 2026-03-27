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
		m := NewGridModel(cols)
		// RenderWidth: (1+10+1) + (1+20+1) = 12 + 22 = 34
		// Outer border: 2
		// Inner divider: 1
		// Total: 34 + 2 + 1 = 37
		assert.Equal(t, 37, m.computeNaturalWidth())
	})

	t.Run("borderless", func(t *testing.T) {
		cols := []Column{
			NewColumn("a", "A", 10),
			NewColumn("b", "B", 20),
		}
		m := NewGridModel(cols).WithBorder(Borderless())
		// RenderWidth: 12 + 22 = 34
		// No borders
		assert.Equal(t, 34, m.computeNaturalWidth())
	})

	t.Run("flex columns use minimum width", func(t *testing.T) {
		cols := []Column{
			NewColumn("a", "A", 10),
			NewFlexColumn("b", "B", 1),
		}
		m := NewGridModel(cols)
		// Fixed: 12
		// Flex minimum: 1 + max(1,1) + 1 = 3 (paddingLeft=1, minWidth=1, paddingRight=1)
		// Outer: 2, Inner: 1
		// Total: 12 + 3 + 2 + 1 = 18
		assert.Equal(t, 18, m.computeNaturalWidth())
	})
}

func TestComputeTotalWidth(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 20),
	}
	m := NewGridModel(cols)

	assert.Equal(t, 37, m.computeTotalWidth())
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
	// Subtract padding: 41 - 1 - 1 - 1 - 1 = 37
	// Flex 1:3, so A gets ~9, B gets ~28
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

func TestFlexColumnMinWidth(t *testing.T) {
	t.Run("default minWidth equals title rune length", func(t *testing.T) {
		col := NewFlexColumn("email", "Email", 1)
		assert.Equal(t, 5, col.MinWidth())
	})

	t.Run("WithMinWidth overrides default", func(t *testing.T) {
		col := NewFlexColumn("email", "Email", 1).WithMinWidth(10)
		assert.Equal(t, 10, col.MinWidth())
	})

	t.Run("fixed columns have zero minWidth", func(t *testing.T) {
		col := NewColumn("name", "Name", 20)
		assert.Equal(t, 0, col.MinWidth())
	})

	t.Run("flex gets minWidth when availableForFlex is zero", func(t *testing.T) {
		cols := []Column{
			NewColumn("a", "A", 40),
			NewFlexColumn("b", "Email", 1),
		}
		border := BorderRounded()
		// totalWidth=47, outer=2, inner=1 → 44 available
		// fixed col render: 1+40+1=42, flex padding: 1+1=2 → used=44, availableForFlex=0
		updateColumnWidths(cols, 47, border)
		assert.Equal(t, 5, cols[1].Width(), "flex should get minWidth (title length)")
	})

	t.Run("flex gets custom minWidth when space is tight", func(t *testing.T) {
		cols := []Column{
			NewColumn("a", "A", 40),
			NewFlexColumn("b", "B", 1).WithMinWidth(8),
		}
		border := BorderRounded()
		// Same tight space scenario
		updateColumnWidths(cols, 45, border)
		assert.Equal(t, 8, cols[1].Width(), "flex should get custom minWidth")
	})

	t.Run("normal flex distribution clamps to minWidth", func(t *testing.T) {
		cols := []Column{
			NewFlexColumn("a", "LongTitle", 10),
			NewFlexColumn("b", "Short", 1),
		}
		border := BorderRounded()
		// totalWidth=20, outer=2, inner=1 → 17 available
		// flex padding: (1+1)+(1+1)=4 → 13 for flex content
		// ratio 10:1, so "Short" gets ~1.18 → 2, but minWidth=5
		updateColumnWidths(cols, 20, border)
		assert.GreaterOrEqual(t, cols[1].Width(), 5, "flex should be at least minWidth")
	})
}

func TestComputeNaturalWidthWithMinWidth(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewFlexColumn("b", "Email", 1),
	}
	m := NewGridModel(cols)
	// Fixed: 1+10+1 = 12
	// Flex minimum: 1 + max(5,1) + 1 = 7 (paddingLeft=1, minWidth=5, paddingRight=1)
	// Outer: 2, Inner: 1
	// Total: 12 + 7 + 2 + 1 = 22
	assert.Equal(t, 22, m.computeNaturalWidth())
}

func TestChromeHeight(t *testing.T) {
	t.Run("full chrome no footer content", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 10)})
		// outer=2, header=1, header_sep=1, footer=0 (no content) = 4
		assert.Equal(t, 4, m.chromeHeight())
	})

	t.Run("full chrome with static footer", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 10)}).
			WithStaticFooter("info")
		// outer=2, header=1, header_sep=1, footer_sep=1, info_row=1 = 6
		assert.Equal(t, 6, m.chromeHeight())
	})

	t.Run("borderless", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 10)}).WithBorder(Borderless())
		// header=1, footer=0 (no content) = 1
		assert.Equal(t, 1, m.chromeHeight())
	})

	t.Run("no header", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 10)}).WithHeaderVisibility(false)
		// outer=2, footer=0 (no content) = 2
		assert.Equal(t, 2, m.chromeHeight())
	})

	t.Run("no footer", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 10)}).WithFooterVisibility(false)
		// outer=2, header=1, header_sep=1 = 4
		assert.Equal(t, 4, m.chromeHeight())
	})
}
