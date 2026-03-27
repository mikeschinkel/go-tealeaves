package teaguide

// ActionSelectedMsg is sent when the user presses an action key in the guide.
// The guide closes itself before emitting this message. The host should
// handle the action identified by ActionKey.
type ActionSelectedMsg struct {
	ActionKey string
}

// GuideDismissedMsg is sent when the user closes the guide without selecting
// an action (via Esc or N).
type GuideDismissedMsg struct{}
