package teagrid

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestHorizontalScrolling(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
		NewColumn("4", "4", 4),
	}).
		WithRows([]Row{
			NewRow(RowData{
				"1": "x1",
				"2": "x2",
				"3": "x3",
				"4": "x4",
			}),
		}).
		WithMaxTotalWidth(18).
		Focused(true)

	const expectedTableOriginal = `┏━━━━┳━━━━┳━━━━┳━┓
┃   1┃   2┃   3┃>┃
┣━━━━╋━━━━╋━━━━╋━┫
┃  x1┃  x2┃  x3┃>┃
┗━━━━┻━━━━┻━━━━┻━┛`

	const expectedTableAfter = `┏━┳━━━━┳━━━━┳━━━━┓
┃<┃   2┃   3┃   4┃
┣━╋━━━━╋━━━━╋━━━━┫
┃<┃  x2┃  x3┃  x4┃
┗━┻━━━━┻━━━━┻━━━━┛`

	hitScrollRight := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftRight})
	}

	hitScrollLeft := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftLeft})
	}

	assert.Equal(t, expectedTableOriginal, model.View())

	hitScrollRight()

	assert.Equal(t, expectedTableAfter, model.View())

	hitScrollLeft()
	assert.Equal(t, expectedTableOriginal, model.View())

	// Try it again, should do nothing
	hitScrollLeft()
	assert.Equal(t, expectedTableOriginal, model.View())
}

func TestHorizontalScrollWithFooter(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
		NewColumn("4", "4", 4),
	}).
		WithRows([]Row{
			NewRow(RowData{
				"1": "x1",
				"2": "x2",
				"3": "x3",
				"4": "x4",
			}),
		}).
		WithStaticFooter("Footer").
		WithMaxTotalWidth(18).
		Focused(true)

	const expectedTableOriginal = `┏━━━━┳━━━━┳━━━━┳━┓
┃   1┃   2┃   3┃>┃
┣━━━━╋━━━━╋━━━━╋━┫
┃  x1┃  x2┃  x3┃>┃
┣━━━━┻━━━━┻━━━━┻━┫
┃          Footer┃
┗━━━━━━━━━━━━━━━━┛`

	const expectedTableAfter = `┏━┳━━━━┳━━━━┳━━━━┓
┃<┃   2┃   3┃   4┃
┣━╋━━━━╋━━━━╋━━━━┫
┃<┃  x2┃  x3┃  x4┃
┣━┻━━━━┻━━━━┻━━━━┫
┃          Footer┃
┗━━━━━━━━━━━━━━━━┛`

	hitScrollRight := func() {
		// Try the programmatic API
		model = model.ScrollRight()
	}

	hitScrollLeft := func() {
		model = model.ScrollLeft()
	}

	assert.Equal(t, expectedTableOriginal, model.View())

	hitScrollRight()

	assert.Equal(t, expectedTableAfter, model.View())

	hitScrollLeft()
	assert.Equal(t, expectedTableOriginal, model.View())

	// Try it again, should do nothing
	hitScrollLeft()
	assert.Equal(t, expectedTableOriginal, model.View())
}

func TestHorizontalScrollingWithFooterAndFrozenCols(t *testing.T) {
	model := New([]Column{
		NewColumn("Name", "Name", 4),
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
		NewColumn("4", "4", 4),
	}).
		WithRows([]Row{
			NewRow(RowData{
				"Name": "A",
				"1":    "x1",
				"2":    "x2",
				"3":    "x3",
				"4":    "x4",
			}),
		}).
		WithStaticFooter("Footer").
		WithMaxTotalWidth(21).
		WithHorizontalFreezeColumnCount(1).
		Focused(true)

	const expectedTableOriginal = `┏━━━━┳━━━━┳━━━━┳━━━━┓
┃Name┃   1┃   2┃   >┃
┣━━━━╋━━━━╋━━━━╋━━━━┫
┃   A┃  x1┃  x2┃   >┃
┣━━━━┻━━━━┻━━━━┻━━━━┫
┃             Footer┃
┗━━━━━━━━━━━━━━━━━━━┛`

	const expectedTableAfter = `┏━━━━┳━┳━━━━┳━━━━┳━━┓
┃Name┃<┃   2┃   3┃ >┃
┣━━━━╋━╋━━━━╋━━━━╋━━┫
┃   A┃<┃  x2┃  x3┃ >┃
┣━━━━┻━┻━━━━┻━━━━┻━━┫
┃             Footer┃
┗━━━━━━━━━━━━━━━━━━━┛`

	hitScrollRight := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftRight})
	}

	hitScrollLeft := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftLeft})
	}

	assert.Equal(t, expectedTableOriginal, model.View())

	hitScrollRight()

	assert.Equal(t, expectedTableAfter, model.View())

	hitScrollLeft()
	assert.Equal(t, expectedTableOriginal, model.View())

	// Try it again, should do nothing
	hitScrollLeft()
	assert.Equal(t, expectedTableOriginal, model.View())
}

// This is long due to literal strings.
func TestHorizontalScrollStopsAtLastColumnBeingVisible(t *testing.T) {
	model := New([]Column{
		NewColumn("Name", "Name", 4),
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
		NewColumn("4", "4", 4),
	}).
		WithRows([]Row{
			NewRow(RowData{
				"Name": "A",
				"1":    "x1",
				"2":    "x2",
				"3":    "x3",
				"4":    "x4",
			}),
		}).
		WithStaticFooter("Footer").
		WithMaxTotalWidth(21).
		WithHorizontalFreezeColumnCount(1).
		Focused(true)

	const expectedTableLeft = `┏━━━━┳━━━━┳━━━━┳━━━━┓
┃Name┃   1┃   2┃   >┃
┣━━━━╋━━━━╋━━━━╋━━━━┫
┃   A┃  x1┃  x2┃   >┃
┣━━━━┻━━━━┻━━━━┻━━━━┫
┃             Footer┃
┗━━━━━━━━━━━━━━━━━━━┛`

	const expectedTableMiddle = `┏━━━━┳━┳━━━━┳━━━━┳━━┓
┃Name┃<┃   2┃   3┃ >┃
┣━━━━╋━╋━━━━╋━━━━╋━━┫
┃   A┃<┃  x2┃  x3┃ >┃
┣━━━━┻━┻━━━━┻━━━━┻━━┫
┃             Footer┃
┗━━━━━━━━━━━━━━━━━━━┛`

	const expectedTableRight = `┏━━━━┳━┳━━━━┳━━━━┓
┃Name┃<┃   3┃   4┃
┣━━━━╋━╋━━━━╋━━━━┫
┃   A┃<┃  x3┃  x4┃
┣━━━━┻━┻━━━━┻━━━━┫
┃          Footer┃
┗━━━━━━━━━━━━━━━━┛`

	hitScrollRight := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftRight})
	}

	assert.Equal(t, expectedTableLeft, model.View())

	hitScrollRight()

	assert.Equal(t, expectedTableMiddle, model.View())

	hitScrollRight()
	assert.Equal(t, expectedTableRight, model.View())

	// Should no longer scroll
	hitScrollRight()
	assert.Equal(t, expectedTableRight, model.View())
}

func TestNoScrollingWhenEntireTableIsVisible(t *testing.T) {
	model := New([]Column{
		NewColumn("Name", "Name", 4),
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	}).
		WithRows([]Row{
			NewRow(RowData{
				"Name": "A",
				"1":    "x1",
				"2":    "x2",
				"3":    "x3",
			}),
		}).
		WithStaticFooter("Footer").
		WithMaxTotalWidth(21).
		WithHorizontalFreezeColumnCount(1).
		Focused(true)

	const expectedTable = `┏━━━━┳━━━━┳━━━━┳━━━━┓
┃Name┃   1┃   2┃   3┃
┣━━━━╋━━━━╋━━━━╋━━━━┫
┃   A┃  x1┃  x2┃  x3┃
┣━━━━┻━━━━┻━━━━┻━━━━┫
┃             Footer┃
┗━━━━━━━━━━━━━━━━━━━┛`

	hitScrollRight := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftRight})
	}

	assert.Equal(t, expectedTable, model.View())

	hitScrollRight()

	assert.Equal(t, expectedTable, model.View())
}

// This is long because of test cases
//
//nolint:funlen
func TestHorizontalScrollingStopEdgeCases(t *testing.T) {
	tests := []struct {
		numCols      int
		nameWidth    int
		colWidth     int
		maxWidth     int
		expectedCols []int
	}{
		{
			numCols:   8,
			nameWidth: 5,
			colWidth:  3,
			maxWidth:  30,
		},
		{
			numCols:      8,
			nameWidth:    5,
			colWidth:     3,
			maxWidth:     20,
			expectedCols: []int{7, 8},
		},
		{
			numCols:   6,
			nameWidth: 5,
			colWidth:  3,
			maxWidth:  30,
		},
		{
			numCols:   50,
			nameWidth: 20,
			colWidth:  6,
			maxWidth:  31,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			cols := []Column{NewColumn("Name", "Name", test.nameWidth)}
			for i := 0; i < test.numCols; i++ {
				s := fmt.Sprintf("%d", i+1)
				cols = append(cols, NewColumn(s, s, test.colWidth))
			}

			rowData := RowData{"Name": "A"}

			for i := 0; i < test.numCols; i++ {
				s := fmt.Sprintf("%d", i+1)
				rowData[s] = s
			}

			rows := []Row{NewRow(rowData)}

			model := New(cols).
				WithRows(rows).
				WithStaticFooter("Footer").
				WithMaxTotalWidth(test.maxWidth).
				WithHorizontalFreezeColumnCount(1).
				Focused(true)

			hitScrollRight := func() {
				model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftRight})
			}

			// Excessive scrolling attempts to be sure
			for i := 0; i < test.numCols*2; i++ {
				hitScrollRight()
			}

			rendered := model.View()

			assert.NotContains(t, rendered, ">")
			assert.Contains(t, rendered, fmt.Sprintf("%d", test.numCols))

			for _, expected := range test.expectedCols {
				assert.Contains(t, rendered, fmt.Sprintf("%d", expected), "Missing expected column")
			}
		})
	}
}

func TestHorizontalScrollingWithCustomKeybind(t *testing.T) {
	keymap := DefaultKeyMap()

	// These intentionally overlap with the keybinds for paging, to ensure
	// that conflicts can live together
	keymap.ScrollRight = key.NewBinding(key.WithKeys("right"))
	keymap.ScrollLeft = key.NewBinding(key.WithKeys("left"))

	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
		NewColumn("4", "4", 4),
	}).
		WithRows([]Row{
			NewRow(RowData{
				"1": "x1",
				"2": "x2",
				"3": "x3",
				"4": "x4",
			}),
		}).
		WithKeyMap(keymap).
		WithMaxTotalWidth(18).
		Focused(true)

	const expectedTableOriginal = `┏━━━━┳━━━━┳━━━━┳━┓
┃   1┃   2┃   3┃>┃
┣━━━━╋━━━━╋━━━━╋━┫
┃  x1┃  x2┃  x3┃>┃
┗━━━━┻━━━━┻━━━━┻━┛`

	const expectedTableAfter = `┏━┳━━━━┳━━━━┳━━━━┓
┃<┃   2┃   3┃   4┃
┣━╋━━━━╋━━━━╋━━━━┫
┃<┃  x2┃  x3┃  x4┃
┗━┻━━━━┻━━━━┻━━━━┛`

	hitScrollRight := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	}

	hitScrollLeft := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	}

	assert.Equal(t, expectedTableOriginal, model.View())

	hitScrollRight()

	assert.Equal(t, expectedTableAfter, model.View())

	hitScrollLeft()
	assert.Equal(t, expectedTableOriginal, model.View())

	// Try it again, should do nothing
	hitScrollLeft()
	assert.Equal(t, expectedTableOriginal, model.View())
}
