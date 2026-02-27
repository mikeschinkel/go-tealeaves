package teastatus

// SetMenuItemsMsg replaces menu items via tea.Msg (for future non-UIViewer use).
type SetMenuItemsMsg struct {
	Items []MenuItem
}

// SetIndicatorsMsg replaces indicators via tea.Msg (for future non-UIViewer use).
type SetIndicatorsMsg struct {
	Indicators []StatusIndicator
}
