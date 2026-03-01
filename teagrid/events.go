package teagrid

// UserEvent is a state change due to user input. These are ONLY
// generated from direct user interaction, not from programmatic changes.
type UserEvent any

// UserEventHighlightedIndexChanged indicates the user scrolled to a new row.
type UserEventHighlightedIndexChanged struct {
	PreviousRowIndex int
	SelectedRowIndex int
}

// UserEventRowSelectToggled indicates a row selection was toggled.
type UserEventRowSelectToggled struct {
	RowIndex   int
	IsSelected bool
}

// UserEventFilterInputFocused indicates the filter input gained focus.
type UserEventFilterInputFocused struct{}

// UserEventFilterInputUnfocused indicates the filter input lost focus.
type UserEventFilterInputUnfocused struct{}

// UserEventCellSelected indicates the user pressed enter/select on a cell.
type UserEventCellSelected struct {
	RowIndex    int
	ColumnIndex int
	ColumnKey   string
	Data        any
}

// --- Editing stubs (implementation deferred to v0.3.0) ---

// CellEditStartedMsg is emitted when cell editing begins.
// Stub: not emitted in v0.2.0.
type CellEditStartedMsg struct {
	RowIndex    int
	ColumnIndex int
	ColumnKey   string
}

// CellEditedMsg is emitted when a cell edit is committed.
// Stub: not emitted in v0.2.0.
type CellEditedMsg struct {
	RowIndex    int
	ColumnIndex int
	ColumnKey   string
	OldValue    any
	NewValue    any
}

// RowEditedMsg is emitted when a row edit is committed.
// Stub: not emitted in v0.2.0.
type RowEditedMsg struct {
	RowIndex int
	OldData  RowData
	NewData  RowData
}

// RowEditCancelledMsg is emitted when a row edit is cancelled.
// Stub: not emitted in v0.2.0.
type RowEditCancelledMsg struct {
	RowIndex int
}

// CellValidatorFunc validates a cell value before committing an edit.
// Stub: accepted but unused in v0.2.0.
type CellValidatorFunc func(columnKey string, value any) error
