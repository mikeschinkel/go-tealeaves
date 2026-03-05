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
func (m GridModel) WithRows(rows []Row) GridModel {
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
func (m GridModel) WithColumns(columns []Column) GridModel {
	m.columns = make([]Column, len(columns))
	copy(m.columns, columns)
	m.recalculateWidth()
	return m
}

// WithKeyMap sets the key bindings.
func (m GridModel) WithKeyMap(keyMap KeyMap) GridModel {
	m.keyMap = keyMap
	return m
}

// WithRowStyleFunc sets a function for dynamic row styling.
// Overrides HighlightStyle when set.
func (m GridModel) WithRowStyleFunc(f func(RowStyleFuncInput) lipgloss.Style) GridModel {
	m.rowStyleFunc = f
	return m
}

// WithHighlightedRow sets the highlighted row index.
func (m GridModel) WithHighlightedRow(index int) GridModel {
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
func (m GridModel) WithCellCursorMode(enabled bool) GridModel {
	m.cellCursorMode = enabled
	return m
}

// WithBaseStyle sets the base style for the grid.
func (m GridModel) WithBaseStyle(style lipgloss.Style) GridModel {
	m.baseStyle = style
	return m
}

// WithHeaderStyle sets the header text style.
func (m GridModel) WithHeaderStyle(style lipgloss.Style) GridModel {
	m.headerStyle = style
	return m
}

// WithHighlightStyle sets the highlighted row style.
func (m GridModel) WithHighlightStyle(style lipgloss.Style) GridModel {
	m.highlightStyle = style
	return m
}

// WithCellCursorStyle sets the cell cursor style.
func (m GridModel) WithCellCursorStyle(style lipgloss.Style) GridModel {
	m.cellCursorStyle = style
	return m
}

// WithFooterStyle sets the footer style (independent from baseStyle).
func (m GridModel) WithFooterStyle(style lipgloss.Style) GridModel {
	m.footerStyle = style
	return m
}

// WithBorder sets the border configuration.
func (m GridModel) WithBorder(border BorderConfig) GridModel {
	m.border = border
	m.recalculateWidth()
	return m
}

// Focused sets whether the grid receives input and shows cursor.
func (m GridModel) Focused(focused bool) GridModel {
	m.focused = focused
	return m
}

// WithSelectableRows enables or disables row selection.
// Unlike v0.1.0, this does NOT auto-add a checkbox column.
// Use WithSelectColumn(true) to add an explicit select column.
func (m GridModel) WithSelectableRows(selectable bool) GridModel {
	m.selectableRows = selectable
	return m
}

// WithSelectColumn adds or removes a visible checkbox column for selection.
// Unlike v0.1.0's SelectableRows, the checkbox column is opt-in.
func (m GridModel) WithSelectColumn(show bool) GridModel {
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
func (m GridModel) WithPageSize(pageSize int) GridModel {
	m.pageSize = pageSize

	maxPages := m.MaxPages()
	if m.currentPage >= maxPages {
		m.currentPage = maxPages - 1
	}

	return m
}

// WithNoPagination disables pagination.
func (m GridModel) WithNoPagination() GridModel {
	m.pageSize = 0
	return m
}

// WithPaginationWrapping sets whether pagination wraps around.
func (m GridModel) WithPaginationWrapping(wrapping bool) GridModel {
	m.paginationWrapping = wrapping
	return m
}

// WithCurrentPage sets the current page (1-indexed).
func (m GridModel) WithCurrentPage(currentPage int) GridModel {
	if m.pageSize == 0 || currentPage == m.CurrentPage() {
		return m
	}

	maxPages := m.MaxPages()
	if currentPage < 1 {
		currentPage = 1
	}
	if currentPage > maxPages {
		currentPage = maxPages
	}

	m.currentPage = currentPage - 1
	m.rowCursorIndex = m.currentPage * m.pageSize
	return m
}

// Filtered enables or disables filtering.
func (m GridModel) Filtered(filtered bool) GridModel {
	m.filtered = filtered
	m.visibleRowCacheUpdated = false
	return m
}

// StartFilterTyping focuses the filter text input.
func (m GridModel) StartFilterTyping() GridModel {
	m.filterTextInput.Focus()
	return m
}

// WithFilterInput sets a custom text input for filtering.
func (m GridModel) WithFilterInput(input textinput.Model) GridModel {
	if m.filterTextInput.Value() != input.Value() {
		m.pageFirst()
	}
	m.filterTextInput = input
	m.visibleRowCacheUpdated = false
	return m
}

// WithFilterInputValue sets the filter string directly.
func (m GridModel) WithFilterInputValue(value string) GridModel {
	if m.filterTextInput.Value() != value {
		m.pageFirst()
	}
	m.filterTextInput.SetValue(value)
	m.filterTextInput.Blur()
	m.visibleRowCacheUpdated = false
	return m
}

// WithFilterFunc sets a custom filter function.
func (m GridModel) WithFilterFunc(fn FilterFunc) GridModel {
	m.filterFunc = fn
	m.visibleRowCacheUpdated = false
	return m
}

// WithFuzzyFilter enables fuzzy (subsequence) filtering.
func (m GridModel) WithFuzzyFilter() GridModel {
	return m.WithFilterFunc(filterFuncFuzzy)
}

// WithStaticFooter sets a static footer text.
func (m GridModel) WithStaticFooter(footer string) GridModel {
	m.staticFooter = footer
	return m
}

// WithHeaderVisibility sets header visibility.
func (m GridModel) WithHeaderVisibility(visible bool) GridModel {
	m.headerVisible = visible
	return m
}

// WithFooterVisibility sets footer visibility.
func (m GridModel) WithFooterVisibility(visible bool) GridModel {
	m.footerVisible = visible
	return m
}

// WithHorizontalFreezeColumnCount freezes columns for horizontal scrolling.
func (m GridModel) WithHorizontalFreezeColumnCount(count int) GridModel {
	m.horizontalScrollFreezeColumnsCount = count
	m.recalculateWidth()
	return m
}

// ScrollRight scrolls one column right.
func (m GridModel) ScrollRight() GridModel {
	m.scrollRight()
	return m
}

// ScrollLeft scrolls one column left.
func (m GridModel) ScrollLeft() GridModel {
	m.scrollLeft()
	return m
}

// PageDown goes to the next page.
func (m GridModel) PageDown() GridModel {
	m.pageDown()
	return m
}

// PageUp goes to the previous page.
func (m GridModel) PageUp() GridModel {
	m.pageUp()
	return m
}

// PageLast goes to the last page.
func (m GridModel) PageLast() GridModel {
	m.pageLast()
	return m
}

// PageFirst goes to the first page.
func (m GridModel) PageFirst() GridModel {
	m.pageFirst()
	return m
}

// WithMissingDataIndicator sets text shown when a row has no data for a column.
func (m GridModel) WithMissingDataIndicator(str string) GridModel {
	m.missingDataIndicator = str
	return m
}

// WithSelectedText sets the display text for selected/unselected states.
func (m GridModel) WithSelectedText(unselected, selected string) GridModel {
	m.selectedText = selected
	m.unselectedText = unselected
	return m
}

// WithAllRowsDeselected clears all row selections.
func (m GridModel) WithAllRowsDeselected() GridModel {
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
func (m GridModel) WithMinimumHeight(height int) GridModel {
	m.minimumHeight = height
	return m
}

// WithMetadata sets grid-level metadata passed to filter and style functions.
func (m GridModel) WithMetadata(metadata map[string]any) GridModel {
	m.metadata = metadata
	return m
}

// WithOverflowIndicator enables or disables overflow indicators for scrolling.
func (m GridModel) WithOverflowIndicator(enabled bool) GridModel {
	m.overflowIndicator = enabled
	return m
}

// --- Editing stubs (v0.3.0) ---

// WithEditable enables or disables cell editing.
// Stub: accepted but non-functional in v0.2.0.
func (m GridModel) WithEditable(editable bool) GridModel {
	m.editable = editable
	return m
}

// WithCellValidator sets a validation function for cell edits.
// Stub: accepted but unused in v0.2.0.
func (m GridModel) WithCellValidator(fn CellValidatorFunc) GridModel {
	m.cellValidator = fn
	return m
}

// --- Help key extensions ---

// WithAdditionalShortHelpKeys adds extra key bindings to short help.
func (m GridModel) WithAdditionalShortHelpKeys(keys []key.Binding) GridModel {
	m.additionalShortHelpKeys = func() []key.Binding { return keys }
	return m
}

// WithAdditionalFullHelpKeys adds extra key bindings to full help.
func (m GridModel) WithAdditionalFullHelpKeys(keys []key.Binding) GridModel {
	m.additionalFullHelpKeys = func() []key.Binding { return keys }
	return m
}
