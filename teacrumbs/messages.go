package teacrumbs

// PushCrumbMsg appends a crumb to the trail via tea.Msg.
type PushCrumbMsg struct {
	Crumb Crumb
}

// PopCrumbMsg removes the last crumb from the trail via tea.Msg.
type PopCrumbMsg struct{}

// SetTrailMsg replaces the entire trail via tea.Msg.
type SetTrailMsg struct {
	Trail []Crumb
}

// SetCrumbMsg updates the crumb at the given index via tea.Msg.
type SetCrumbMsg struct {
	Index int
	Crumb Crumb
}
