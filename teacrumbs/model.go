package teacrumbs

import (
	tea "charm.land/bubbletea/v2"
)

// BreadcrumbsModel is a Bubble Tea model for a breadcrumb trail.
// It manages the trail of crumbs and renders them with styling and truncation.
type BreadcrumbsModel struct {
	Styles    Styles
	trail     []Crumb
	separator string
	width     int
}

// NewBreadcrumbsModel creates a new BreadcrumbsModel with default styles.
func NewBreadcrumbsModel() BreadcrumbsModel {
	return BreadcrumbsModel{
		Styles:    DefaultStyles(),
		separator: " > ",
	}
}

// Push appends a crumb to the trail and returns a copy.
func (m BreadcrumbsModel) Push(crumb Crumb) BreadcrumbsModel {
	newTrail := make([]Crumb, len(m.trail)+1)
	copy(newTrail, m.trail)
	newTrail[len(m.trail)] = crumb
	m.trail = newTrail
	return m
}

// Pop removes the last crumb from the trail and returns a copy.
// No-op if the trail is empty.
func (m BreadcrumbsModel) Pop() BreadcrumbsModel {
	if len(m.trail) == 0 {
		return m
	}
	newTrail := make([]Crumb, len(m.trail)-1)
	copy(newTrail, m.trail[:len(m.trail)-1])
	m.trail = newTrail
	return m
}

// SetTrail replaces the entire trail and returns a copy.
func (m BreadcrumbsModel) SetTrail(trail []Crumb) BreadcrumbsModel {
	newTrail := make([]Crumb, len(trail))
	copy(newTrail, trail)
	m.trail = newTrail
	return m
}

// SetCrumb updates the crumb at the given index and returns a copy.
// No-op if the index is out of range.
func (m BreadcrumbsModel) SetCrumb(index int, crumb Crumb) BreadcrumbsModel {
	if index < 0 || index >= len(m.trail) {
		return m
	}
	newTrail := make([]Crumb, len(m.trail))
	copy(newTrail, m.trail)
	newTrail[index] = crumb
	m.trail = newTrail
	return m
}

// Trail returns a copy of the current trail.
func (m BreadcrumbsModel) Trail() []Crumb {
	out := make([]Crumb, len(m.trail))
	copy(out, m.trail)
	return out
}

// Len returns the number of crumbs in the trail.
func (m BreadcrumbsModel) Len() int {
	return len(m.trail)
}

// SetSize sets the available width for the breadcrumb trail.
func (m BreadcrumbsModel) SetSize(width int) BreadcrumbsModel {
	m.width = width
	return m
}

// WithStyles returns a copy with the given styles override.
func (m BreadcrumbsModel) WithStyles(styles Styles) BreadcrumbsModel {
	m.Styles = styles
	return m
}

// WithSeparator returns a copy with a custom separator string.
func (m BreadcrumbsModel) WithSeparator(sep string) BreadcrumbsModel {
	m.separator = sep
	return m
}

// Init implements tea.Model. No-op for breadcrumbs.
func (m BreadcrumbsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model. Handles breadcrumb messages and WindowSizeMsg.
func (m BreadcrumbsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case PushCrumbMsg:
		m = m.Push(msg.Crumb)
	case PopCrumbMsg:
		m = m.Pop()
	case SetTrailMsg:
		m = m.SetTrail(msg.Trail)
	case SetCrumbMsg:
		m = m.SetCrumb(msg.Index, msg.Crumb)
	}
	return m, nil
}

// View implements tea.Model. Renders the breadcrumb trail.
func (m BreadcrumbsModel) View() tea.View {
	return tea.NewView(renderTrail(m.trail, m.separator, m.width, m.Styles))
}
