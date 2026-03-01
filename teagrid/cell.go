package teagrid

import "charm.land/lipgloss/v2"

// CellValue represents a value in a grid cell with optional styling,
// sort key, and rich text content. It replaces the v0.1.0 StyledCell type.
type CellValue struct {
	// Data is the display value for the cell.
	Data any

	// SortValue, if non-nil, is used for sorting instead of Data.
	// This allows displaying a human-friendly value while sorting
	// by a machine-friendly key (e.g., display "Jan 1" but sort by timestamp).
	SortValue any

	// Style is the cell-level style. Ignored if StyleFunc is set.
	Style lipgloss.Style

	// StyleFunc returns a dynamic style based on cell context.
	// When set, overrides Style.
	StyleFunc CellStyleFunc

	// Spans, if non-empty, override Data for display purposes.
	// Each span can have its own style for rich inline text.
	Spans []Span
}

// Span represents a styled segment of text within a cell.
type Span struct {
	Text  string
	Style lipgloss.Style
}

// CellStyleInput provides context to CellStyleFunc for dynamic styling.
type CellStyleInput struct {
	// Data is the cell's data value.
	Data any

	// Column is the column this cell belongs to.
	Column Column

	// Row is the row this cell belongs to.
	Row Row

	// RowIndex is the index of the row in the visible rows.
	RowIndex int

	// ColumnIndex is the index of the column.
	ColumnIndex int

	// IsHighlightedRow is true when this cell's row has the row cursor.
	IsHighlightedRow bool

	// IsCursorCell is true when this specific cell has the cell cursor.
	IsCursorCell bool

	// GlobalMetadata is the grid-level metadata set via WithMetadata.
	GlobalMetadata map[string]any
}

// CellStyleFunc returns a style based on cell context.
// It receives cursor and highlight state so consumers can customize
// appearance without rebuilding rows.
type CellStyleFunc func(CellStyleInput) lipgloss.Style

// NewCellValue creates a CellValue with the given data and style.
func NewCellValue(data any, style lipgloss.Style) CellValue {
	return CellValue{
		Data:  data,
		Style: style,
	}
}

// NewCellValueWithStyleFunc creates a CellValue with a dynamic style function.
func NewCellValueWithStyleFunc(data any, fn CellStyleFunc) CellValue {
	return CellValue{
		Data:      data,
		StyleFunc: fn,
	}
}

// NewCellValueWithSortKey creates a CellValue with separate display and sort values.
func NewCellValueWithSortKey(data any, sortValue any, style lipgloss.Style) CellValue {
	return CellValue{
		Data:      data,
		SortValue: sortValue,
		Style:     style,
	}
}

// NewCellValueWithSpans creates a CellValue with rich text spans.
func NewCellValueWithSpans(spans []Span, style lipgloss.Style) CellValue {
	return CellValue{
		Spans: spans,
		Style: style,
	}
}

// NewSpan creates a styled text span.
func NewSpan(text string, style lipgloss.Style) Span {
	return Span{
		Text:  text,
		Style: style,
	}
}

// SortableValue returns the value to use for sorting.
// If SortValue is set, returns SortValue; otherwise Data.
func (cv CellValue) SortableValue() any {
	if cv.SortValue != nil {
		return cv.SortValue
	}
	return cv.Data
}

// HasSpans returns whether the cell has rich text spans.
func (cv CellValue) HasSpans() bool {
	return len(cv.Spans) > 0
}
