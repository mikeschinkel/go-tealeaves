package teagrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventTypes(t *testing.T) {
	// Verify all event types can be used as UserEvent
	var events []UserEvent

	events = append(events,
		UserEventHighlightedIndexChanged{PreviousRowIndex: 0, SelectedRowIndex: 1},
		UserEventRowSelectToggled{RowIndex: 0, IsSelected: true},
		UserEventFilterInputFocused{},
		UserEventFilterInputUnfocused{},
		UserEventCellSelected{RowIndex: 0, ColumnIndex: 1, ColumnKey: "name", Data: "test"},
	)

	assert.Len(t, events, 5)
}

func TestUserEventCellSelected(t *testing.T) {
	event := UserEventCellSelected{
		RowIndex:    2,
		ColumnIndex: 1,
		ColumnKey:   "name",
		Data:        "Alice",
	}

	assert.Equal(t, 2, event.RowIndex)
	assert.Equal(t, 1, event.ColumnIndex)
	assert.Equal(t, "name", event.ColumnKey)
	assert.Equal(t, "Alice", event.Data)
}

func TestEditingStubTypes(t *testing.T) {
	// Verify editing stub types compile and can be instantiated
	_ = CellEditStartedMsg{RowIndex: 0, ColumnIndex: 1, ColumnKey: "x"}
	_ = CellEditedMsg{RowIndex: 0, ColumnIndex: 1, ColumnKey: "x", OldValue: "a", NewValue: "b"}
	_ = RowEditedMsg{RowIndex: 0, OldData: RowData{"x": 1}, NewData: RowData{"x": 2}}
	_ = RowEditCancelledMsg{RowIndex: 0}

	// CellValidatorFunc should be assignable
	var validator CellValidatorFunc = func(columnKey string, value any) error {
		return nil
	}
	assert.Nil(t, validator("x", "value"))
}
