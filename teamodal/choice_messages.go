package teamodal

// ChoiceSelectedMsg is sent when a choice is selected (Enter on focused button or hotkey)
type ChoiceSelectedMsg struct {
	OptionID string // The ID of the selected option
	Index    int    // The 0-based index of the selected option
}

// ChoiceCancelledMsg is sent when the modal is cancelled (Esc key)
type ChoiceCancelledMsg struct{}
