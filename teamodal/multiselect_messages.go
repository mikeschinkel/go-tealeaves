package teamodal

// MultiSelectButtonPressedMsg is sent when user presses a button (any button).
// All buttons emit this message; the consumer maps ButtonID to behavior.
// Esc is the only path to MultiSelectCancelledMsg.
type MultiSelectButtonPressedMsg[T MultiSelectItem] struct {
	ButtonID string // Which button was pressed
	Selected []T    // Items that were checked at the time
}

// MultiSelectCancelledMsg is sent when user presses Esc
type MultiSelectCancelledMsg struct{}
