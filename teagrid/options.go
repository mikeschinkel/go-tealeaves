package teagrid

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

// RowStyleFuncInput provides context for row-level dynamic styling.
type RowStyleFuncInput struct {
	Index         int
	Row           Row
	IsHighlighted bool
}

// WithRows sets the data rows.
func (m Model) WithRows(rows []Row) Model {
	m.rows = rows
	m.visibleRowCacheUpdated = false

	if m.rowCursorIndex >= len(m.rows) {
		m.rowCursorIndex = len(m.rows) - 1
	}
	if m.rowCursorIndex < 0 {
		m.rowCursorIndex = 0
	}

	if m.pageSize != 0 {
		maxPage := m.MaxPages()
		if maxPage <= m.currentPage {
			m.pageLast()
		}
	}

	return m
}

// WithColumns sets the visible columns.
func (m Model) WithColumns(columns []Column) Model {
	m.columns = make([]Column, len(columns))
	copy(m.columns, columns)
	m.recalculateWidth()
	return m
}

// WithKeyMap sets the key bindings.
func (m Model) WithKeyMap(keyMap KeyMap) Model {
	m.keyMap = keyMap
	return m
}

// WithRowStyleFunc sets a function for dynamic row styling.
// Overrides HighlightStyle when set.
func (m Model) WithRowStyleFunc(f func(RowStyleFuncInput) lipgloss.Style) Model {
	m.rowStyleFunc = f
	return m
}

// WithHighlightedRow sets the highlighted row index.
func (m Model) WithHighlightedRow(index int) Model {
	m.rowCursorIndex = index

	numRows := len(m.GetVisibleRows())
	if m.rowCursorIndex >= numRows {
		m.rowCursorIndex = numRows - 1
	}
	if m.rowCursorIndex < 0 {
		m.rowCursorIndex = 0
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	return m
}

// WithCellCursorMode enables or disables cell cursor navigation.
func (m Model) WithCellCursorMode(enabled bool) Model {
	m.cellCursorMode = enabled
	return m
}

// WithBaseStyle sets the base style for the grid.
func (m Model) WithBaseStyle(style lipgloss.Style) Model {
	m.baseStyle = style
	return m
}

// WithHeaderStyle sets the header text style.
func (m Model) WithHeaderStyle(style lipgloss.Style) Model {
	m.headerStyle = style
	return m
}

// WithHighlightStyle sets the highlighted row style.
func (m Model) WithHighlightStyle(style lipgloss.Style) Model {
	m.highlightStyle = style
	return m
}

// WithCellCursorStyle sets the cell cursor style.
func (m Model) WithCellCursorStyle(style lipgloss.Style) Model {
	m.cellCursorStyle = style
	return m
}

// WithFooterStyle sets the footer style (independent from baseStyle).
func (m Model) WithFooterStyle(style lipgloss.Style) Model {
	m.footerStyle = style
	return m
}

// WithBorder sets the border configuration.
func (m Model) WithBorder(border BorderConfig) Model {
	m.border = border
	m.recalculateWidth()
	return m
}

// Focused sets whether the grid receives input and shows cursor.
func (m Model) Focused(focused bool) Model {
	m.focused = focused
	return m
}

// WithSelectableRows enables or disables row selection.
// Unlike v0.1.0, this does NOT auto-add a checkbox column.
// Use WithSelectColumn(true) to add an explicit select column.
func (m Model) WithSelectableRows(selectable bool) Model {
	m.selectableRows = selectable
	return m
}

// WithSelectColumn adds or removes a visible checkbox column for selection.
// Unlike v0.1.0's SelectableRows, the checkbox column is opt-in.
func (m Model) WithSelectColumn(show bool) Model {
	hadSelectColumn := m.selectColumn && len(m.columns) > 0 &&
		m.columns[0].key == selectColumnKey
	m.selectColumn = show

	if show && !hadSelectColumn {
		selectCol := NewColumn(selectColumnKey, m.selectedText, len([]rune(m.selectedText))).
			WithPaddingLeft(0).WithPaddingRight(0)
		m.columns = append([]Column{selectCol}, m.columns...)
		m.recalculateWidth()
	} else if !show && hadSelectColumn {
		m.columns = m.columns[1:]
		m.recalculateWidth()
	}

	return m
}

// WithPageSize sets the page size for pagination.
func (m Model) WithPageSize(pageSize int) Model {
	m.pageSize = pageSize

	maxPages := m.MaxPages()
	if m.currentPage >= maxPages {
		m.currentPage = maxPages - 1
	}

	return m
}

// WithNoPagination disables pagination.
func (m Model) WithNoPagination() Model {
	m.pageSize = 0
	return m
}

// WithPaginationWrapping sets whether pagination wraps around.
func (m Model) WithPaginationWrapping(wrapping bool) Model {
	m.paginationWrapping = wrapping
	return m
}

// WithCurrentPage sets the current page (1-indexed).
func (m Model) WithCurrentPage(currentPage int) Model {
	if m.pageSize == 0 || currentPage == m.CurrentPage() {
		return m
	}

	if currentPage < 1 {
		currentPage = 1
	} else if maxPages := m.MaxPages(); currentPage > maxPages {
		currentPage = maxPages
	}

	m.currentPage = currentPage - 1
	m.rowCursorIndex = m.currentPage * m.pageSize
	return m
}

// Filtered enables or disables filtering.
func (m Model) Filtered(filtered bool) Model {
	m.filtered = filtered
	m.visibleRowCacheUpdated = false
	return m
}

// StartFilterTyping focuses the filter text input.
func (m Model) StartFilterTyping() Model {
	m.filterTextInput.Focus()
	return m
}

// WithFilterInput sets a custom text input for filtering.
func (m Model) WithFilterInput(input textinput.Model) Model {
	if m.filterTextInput.Value() != input.Value() {
		m.pageFirst()
	}
	m.filterTextInput = input
	m.visibleRowCacheUpdated = false
	return m
}

// WithFilterInputValue sets the filter string directly.
func (m Model) WithFilterInputValue(value string) Model {
	if m.filterTextInput.Value() != value {
		m.pageFirst()
	}
	m.filterTextInput.SetValue(value)
	m.filterTextInput.Blur()
	m.visibleRowCacheUpdated = false
	return m
}

// WithFilterFunc sets a custom filter function.
func (m Model) WithFilterFunc(fn FilterFunc) Model {
	m.filterFunc = fn
	m.visibleRowCacheUpdated = false
	return m
}

// WithFuzzyFilter enables fuzzy (subsequence) filtering.
func (m Model) WithFuzzyFilter() Model {
	return m.WithFilterFunc(filterFuncFuzzy)
}

// WithStaticFooter sets a static footer text.
func (m Model) WithStaticFooter(footer string) Model {
	m.staticFooter = footer
	return m
}

// WithHeaderVisibility sets header visibility.
func (m Model) WithHeaderVisibility(visible bool) Model {
	m.headerVisible = visible
	return m
}

// WithFooterVisibility sets footer visibility.
func (m Model) WithFooterVisibility(visible bool) Model {
	m.footerVisible = visible
	return m
}

// WithHorizontalFreezeColumnCount freezes columns for horizontal scrolling.
func (m Model) WithHorizontalFreezeColumnCount(count int) Model {
	m.horizontalScrollFreezeColumnsCount = count
	m.recalculateWidth()
	return m
}

// ScrollRight scrolls one column right.
func (m Model) ScrollRight() Model {
	m.scrollRight()
	return m
}

// ScrollLeft scrolls one column left.
func (m Model) ScrollLeft() Model {
	m.scrollLeft()
	return m
}

// PageDown goes to the next page.
func (m Model) PageDown() Model {
	m.pageDown()
	return m
}

// PageUp goes to the previous page.
func (m Model) PageUp() Model {
	m.pageUp()
	return m
}

// PageLast goes to the last page.
func (m Model) PageLast() Model {
	m.pageLast()
	return m
}

// PageFirst goes to the first page.
func (m Model) PageFirst() Model {
	m.pageFirst()
	return m
}

// WithMissingDataIndicator sets text shown when a row has no data for a column.
func (m Model) WithMissingDataIndicator(str string) Model {
	m.missingDataIndicator = str
	return m
}

// WithSelectedText sets the display text for selected/unselected states.
func (m Model) WithSelectedText(unselected, selected string) Model {
	m.selectedText = selected
	m.unselectedText = unselected
	return m
}

// WithAllRowsDeselected clears all row selections.
func (m Model) WithAllRowsDeselected() Model {
	rows := m.GetVisibleRows()
	for i, row := range rows {
		if row.selected {
			rows[i] = row.Selected(false)
		}
	}
	m.rows = rows
	return m
}

// WithMinimumHeight sets the minimum total height including borders.
func (m Model) WithMinimumHeight(height int) Model {
	m.minimumHeight = height
	return m
}

// WithMetadata sets grid-level metadata passed to filter and style functions.
func (m Model) WithMetadata(metadata map[string]any) Model {
	m.metadata = metadata
	return m
}

// WithOverflowIndicator enables or disables overflow indicators for scrolling.
func (m Model) WithOverflowIndicator(enabled bool) Model {
	m.overflowIndicator = enabled
	return m
}

// --- Editing stubs (v0.3.0) ---

// WithEditable enables or disables cell editing.
// Stub: accepted but non-functional in v0.2.0.
func (m Model) WithEditable(editable bool) Model {
	m.editable = editable
	return m
}

// WithCellValidator sets a validation function for cell edits.
// Stub: accepted but unused in v0.2.0.
func (m Model) WithCellValidator(fn CellValidatorFunc) Model {
	m.cellValidator = fn
	return m
}

// --- Help key extensions ---

// WithAdditionalShortHelpKeys adds extra key bindings to short help.
func (m Model) WithAdditionalShortHelpKeys(keys []key.Binding) Model {
	m.additionalShortHelpKeys = func() []key.Binding { return keys }
	return m
}

// WithAdditionalFullHelpKeys adds extra key bindings to full help.
func (m Model) WithAdditionalFullHelpKeys(keys []key.Binding) Model {
	m.additionalFullHelpKeys = func() []key.Binding { return keys }
	return m
}
