package teagrid

import (
	"sync/atomic"

	"charm.land/lipgloss/v2"
)

// RowData is a map of column keys to cell values.
// Values can be any type: plain values are displayed with fmt.Sprintf,
// CellValue instances provide styling, sort keys, and rich text.
// Data with keys that don't match any column is retained but not displayed.
type RowData map[string]any

// Row represents a row in the grid.
type Row struct {
	Style lipgloss.Style
	Data  RowData

	selected bool

	// id is an internal unique ID to match rows after they're copied.
	id uint32
}

var lastRowID atomic.Uint32

func init() {
	lastRowID.Store(1)
}

// NewRow creates a new row with the given data.
// Data is shallow-copied to prevent external mutation.
func NewRow(data RowData) Row {
	row := Row{
		Data: make(RowData, len(data)),
		id:   lastRowID.Add(1),
	}

	for k, v := range data {
		row.Data[k] = v
	}

	return row
}

// WithStyle sets the row's style.
func (r Row) WithStyle(style lipgloss.Style) Row {
	r.Style = style
	return r
}

// Selected returns a copy of the row with the given selection state.
func (r Row) Selected(selected bool) Row {
	r.selected = selected
	return r
}

// IsSelected returns whether the row is selected.
func (r Row) IsSelected() bool {
	return r.selected
}

// ID returns the row's internal unique identifier.
func (r Row) ID() uint32 {
	return r.id
}
