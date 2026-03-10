package teacrumbs

import tea "charm.land/bubbletea/v2"

// PushCrumbMsg appends a crumb to the crumbs via tea.Msg.
type PushCrumbMsg struct {
	Crumb Crumb
}

// PopCrumbMsg removes the last crumb from the crumbs via tea.Msg.
type PopCrumbMsg struct{}

// SetCrumbsMsg replaces the entire crumbs via tea.Msg.
type SetCrumbsMsg struct {
	Crumbs []Crumb
}

// SetCrumbMsg updates the crumb at the given index via tea.Msg.
type SetCrumbMsg struct {
	Index int
	Crumb Crumb
}

// CrumbClickedMsg is emitted when a breadcrumb is clicked.
type CrumbClickedMsg struct {
	Index  int
	Crumb  Crumb
	Button tea.MouseButton
}

// CrumbHoverMsg is emitted when the mouse enters a breadcrumb.
type CrumbHoverMsg struct {
	Index int
	Crumb Crumb
}

// CrumbHoverLeaveMsg is emitted when the mouse leaves all breadcrumbs.
type CrumbHoverLeaveMsg struct{}
