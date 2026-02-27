package teamodal

// AnsweredYesMsg is sent when user confirms with Enter on Yes button
type AnsweredYesMsg struct{}

// AnsweredNoMsg is sent when user selects No button and presses Enter
type AnsweredNoMsg struct{}

// ClosedMsg is sent when user closes an alert (OK button) or cancels with Esc
type ClosedMsg struct{}

// ProgressCancelledMsg is sent when user cancels a progress modal with Esc
type ProgressCancelledMsg struct{}

// ProgressBackgroundMsg is sent when user sends a progress operation to background
type ProgressBackgroundMsg struct{}
