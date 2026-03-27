package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestScrollRight(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}
	m := NewGridModel(cols).WithSize(30, 24) // narrow viewport forces scrolling

	assert.Equal(t, 0, m.horizontalScrollOffsetCol)

	m = m.scrollRight()
	if m.maxHorizontalColumnIndex > 0 {
		assert.Equal(t, 1, m.horizontalScrollOffsetCol)
	}
}

func TestScrollLeft(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("a", "A", 20)})
	m.horizontalScrollOffsetCol = 2

	m = m.scrollLeft()
	assert.Equal(t, 1, m.horizontalScrollOffsetCol)

	m = m.scrollLeft()
	assert.Equal(t, 0, m.horizontalScrollOffsetCol)

	m = m.scrollLeft()
	assert.Equal(t, 0, m.horizontalScrollOffsetCol, "should not go below 0")
}

func TestScrollNoOverflowNoScroll(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
	}
	m := NewGridModel(cols).WithSize(80, 24)

	assert.Equal(t, 0, m.maxHorizontalColumnIndex)
}

func TestVisibleColumns(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 10),
		NewColumn("c", "C", 10),
	}

	t.Run("all fit", func(t *testing.T) {
		m := NewGridModel(cols).WithSize(100, 24)
		visible := m.visibleColumns()
		assert.Len(t, visible, 3)
	})

	t.Run("no viewport returns all", func(t *testing.T) {
		m := NewGridModel(cols)
		visible := m.visibleColumns()
		assert.Len(t, visible, 3)
	})
}

func TestHorizontalFreezeColumns(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 5),
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}

	m := NewGridModel(cols).
		WithHorizontalFreezeColumnCount(1).
		WithSize(30, 24)

	// First column should always be in visible set
	visible := m.visibleColumns()
	assert.Equal(t, "id", visible[0].Key())
}

func TestVisibleColumnsNarrowViewport(t *testing.T) {
	// Each column: 10 content + 1 paddingLeft + 1 paddingRight = 12 render width
	// 3 columns + 2 inner dividers + 2 outer borders = 12*3 + 2 + 2 = 40
	// Set viewport to 26 — should fit ~2 columns max
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 10),
		NewColumn("c", "C", 10),
	}

	m := NewGridModel(cols).WithSize(26, 24)
	visible := m.visibleColumns()
	assert.Less(t, len(visible), 3, "narrow viewport should clip columns")
	assert.Greater(t, len(visible), 0, "should show at least one column")
}

func TestRenderOutputWidthRespectViewport(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 10),
		NewColumn("c", "C", 10),
		NewColumn("d", "D", 10),
	}

	viewportWidth := 30
	m := NewGridModel(cols).WithSize(viewportWidth, 24).
		WithRows([]Row{
			NewRow(RowData{"a": "hello", "b": "world", "c": "foo", "d": "bar"}),
		})

	rendered := m.render()
	for i, line := range splitLines(rendered) {
		lineWidth := lipgloss.Width(line)
		assert.LessOrEqual(t, lineWidth, viewportWidth,
			"line %d width %d exceeds viewport %d", i, lineWidth, viewportWidth)
	}
}

func TestColCursorAutoScrollRight(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 10),
		NewColumn("c", "C", 10),
		NewColumn("d", "D", 10),
	}

	// Narrow viewport: only ~2 columns fit
	m := NewGridModel(cols).
		WithSize(26, 24).
		WithColCursorMode(true).
		WithFocused(true)

	// Move cell cursor right repeatedly
	for i := 0; i < 3; i++ {
		m = m.moveColRight()
	}

	// Cursor should be at column 3 ("d")
	assert.Equal(t, 3, m.colCursorColumnIndex)

	// The cursor column should be in the visible set
	visible := m.visibleColumns()
	found := false
	for _, col := range visible {
		if col.Key() == "d" {
			found = true
			break
		}
	}
	assert.True(t, found, "cursor column 'd' should be visible after auto-scroll")
}

func TestColCursorAutoScrollLeft(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 10),
		NewColumn("c", "C", 10),
		NewColumn("d", "D", 10),
	}

	// Start scrolled right
	m := NewGridModel(cols).
		WithSize(26, 24).
		WithColCursorMode(true).
		WithFocused(true)

	// Move right to force scroll, then move back left
	for i := 0; i < 3; i++ {
		m = m.moveColRight()
	}
	// Now move left back to column "a"
	for i := 0; i < 3; i++ {
		m = m.moveColLeft()
	}

	assert.Equal(t, 0, m.colCursorColumnIndex)

	visible := m.visibleColumns()
	assert.Equal(t, "a", visible[0].Key(), "column 'a' should be visible after scrolling left")
}

func TestFreezeDividerInOutput(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 5),
		NewColumn("name", "Name", 10),
		NewColumn("value", "Value", 10),
	}
	rows := []Row{
		NewRow(RowData{"id": "1", "name": "Alice", "value": "100"}),
	}

	m := NewGridModel(cols).
		WithHorizontalFreezeColumnCount(1).
		WithRows(rows).
		WithSize(80, 24)

	rendered := m.render()
	lines := splitLines(rendered)

	// lines[1] is the header row (after top border)
	assert.Contains(t, lines[1], "║",
		"header row should contain freeze divider ║")

	// lines[3] is the first data row (after header separator)
	assert.Contains(t, lines[3], "║",
		"data row should contain freeze divider ║")
}

func TestFreezeDividerJunctions(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 5),
		NewColumn("name", "Name", 10),
		NewColumn("value", "Value", 10),
	}
	rows := []Row{
		NewRow(RowData{"id": "1", "name": "Alice", "value": "100"}),
	}

	// Disable footer so bottom border sits directly below data rows
	m := NewGridModel(cols).
		WithHorizontalFreezeColumnCount(1).
		WithRows(rows).
		WithFooterVisibility(false).
		WithSize(80, 24)

	rendered := m.render()
	lines := splitLines(rendered)

	// Top border should have ╥ at freeze boundary
	assert.Contains(t, lines[0], "╥",
		"top border should contain freeze top junction ╥")

	// Header separator should have ╫ at freeze boundary
	assert.Contains(t, lines[2], "╫",
		"header separator should contain freeze inner junction ╫")

	// Bottom border (last line) should have ╨ at freeze boundary
	assert.Contains(t, lines[len(lines)-1], "╨",
		"bottom border should contain freeze bottom junction ╨")
}

func TestFreezeDividerBottomBorderPlainWithFooter(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 5),
		NewColumn("name", "Name", 10),
		NewColumn("value", "Value", 10),
	}
	rows := []Row{
		NewRow(RowData{"id": "1", "name": "Alice", "value": "100"}),
	}

	// With footer (pagination), bottom border should be plain (no orphaned junctions)
	m := NewGridModel(cols).
		WithHorizontalFreezeColumnCount(1).
		WithRows(rows).
		WithSize(80, 24)

	rendered := m.render()
	lines := splitLines(rendered)

	// Bottom border should NOT have freeze junction — info row above has no column dividers
	assert.NotContains(t, lines[len(lines)-1], "╨",
		"bottom border should NOT contain freeze junction when footer info row is above")
}

func TestNoFreezeDividerWithoutFreezeColumns(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 5),
		NewColumn("name", "Name", 10),
	}
	rows := []Row{
		NewRow(RowData{"id": "1", "name": "Alice"}),
	}

	m := NewGridModel(cols).
		WithRows(rows).
		WithSize(80, 24)

	rendered := m.render()

	assert.NotContains(t, rendered, "║", "should not contain ║ without freeze columns")
	assert.NotContains(t, rendered, "╥", "should not contain ╥ without freeze columns")
	assert.NotContains(t, rendered, "╨", "should not contain ╨ without freeze columns")
	assert.NotContains(t, rendered, "╫", "should not contain ╫ without freeze columns")
}

func TestHasHiddenColumnsLeft(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}
	m := NewGridModel(cols).WithSize(30, 24)

	assert.False(t, m.hasHiddenColumnsLeft(), "no hidden left at offset 0")

	m.horizontalScrollOffsetCol = 1
	assert.True(t, m.hasHiddenColumnsLeft(), "hidden left at offset 1")
}

func TestHasHiddenColumnsRight(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}
	m := NewGridModel(cols).WithSize(30, 24)

	if m.maxHorizontalColumnIndex > 0 {
		assert.True(t, m.hasHiddenColumnsRight(), "hidden right at offset 0 with overflow")

		m.horizontalScrollOffsetCol = m.maxHorizontalColumnIndex
		assert.False(t, m.hasHiddenColumnsRight(), "no hidden right at max offset")
	}
}

func TestLeftVertical(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}

	t.Run("no freeze columns", func(t *testing.T) {
		m := NewGridModel(cols).
			WithSize(30, 24).
			WithOverflowIndicator(true)

		// At offset 0: no hidden left → normal vertical
		assert.Equal(t, m.border.Chars.Vertical, m.leftVertical())

		// Scroll right: hidden left → overflow vertical on outer border
		m.horizontalScrollOffsetCol = 1
		assert.Equal(t, m.border.Chars.OverflowVertical, m.leftVertical())
	})

	t.Run("with freeze columns uses normal outer border", func(t *testing.T) {
		m := NewGridModel(cols).
			WithSize(30, 24).
			WithOverflowIndicator(true).
			WithHorizontalFreezeColumnCount(1)

		// Even with hidden left, outer border stays normal when freeze > 0
		m.horizontalScrollOffsetCol = 1
		assert.Equal(t, m.border.Chars.Vertical, m.leftVertical(),
			"left outer border should be normal when freeze columns handle overflow")
	})
}

func TestRightVertical(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}
	m := NewGridModel(cols).
		WithSize(30, 24).
		WithOverflowIndicator(true)

	if m.maxHorizontalColumnIndex > 0 {
		// At offset 0 with overflow: hidden right → overflow vertical
		assert.Equal(t, m.border.Chars.OverflowVertical, m.rightVertical())

		// At max offset: no hidden right → normal vertical
		m.horizontalScrollOffsetCol = m.maxHorizontalColumnIndex
		assert.Equal(t, m.border.Chars.Vertical, m.rightVertical())
	}
}

func TestFillWidthEnabledByDefault(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)})
	assert.True(t, m.fillWidth, "fillWidth should be true by default")
}

func TestFillWidthMatchesViewport(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 10),
		NewColumn("c", "C", 10),
		NewColumn("d", "D", 10),
	}
	rows := []Row{
		NewRow(RowData{"a": "hello", "b": "world", "c": "foo", "d": "bar"}),
	}

	viewportWidth := 30
	m := NewGridModel(cols).
		WithFillWidth(true).
		WithRows(rows).
		WithSize(viewportWidth, 24)

	// Every rendered line should match the viewport width exactly
	rendered := m.render()
	for i, line := range splitLines(rendered) {
		lineWidth := lipgloss.Width(line)
		assert.Equal(t, viewportWidth, lineWidth,
			"line %d width %d should equal viewport %d", i, lineWidth, viewportWidth)
	}
}

func TestFillWidthDoesNotMutateOriginalColumns(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 10),
		NewColumn("c", "C", 10),
		NewColumn("d", "D", 10),
	}

	m := NewGridModel(cols).
		WithFillWidth(true).
		WithSize(30, 24)

	// Force visible column computation
	_ = m.visibleColumns()

	// Original columns should be untouched
	for i, col := range m.columns {
		assert.Equal(t, defaultPaddingRight, col.paddingRight,
			"column %d paddingRight should be unchanged", i)
	}
}

func TestOverflowVerticalEnabledByDefault(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 20),
		NewColumn("b", "B", 20),
		NewColumn("c", "C", 20),
	}
	// overflowIndicator is true by default
	m := NewGridModel(cols).WithSize(30, 24)

	m.horizontalScrollOffsetCol = 1
	assert.Equal(t, m.border.Chars.OverflowVertical, m.leftVertical(),
		"should use overflow vertical when overflow indicator is enabled by default")
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
