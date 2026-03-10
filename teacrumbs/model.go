package teacrumbs

import (
	tea "charm.land/bubbletea/v2"
)

// mouseState holds mutable mouse tracking state shared via pointer.
// This survives value copies of BreadcrumbsModel (same pattern as teamodal).
type mouseState struct {
	crumbBounds []crumbBound
	hoveredIdx  int
}

// crumbBound tracks the horizontal extent of a rendered crumb.
type crumbBound struct {
	startX, endX int // Relative to crumbs start
}

// BreadcrumbsModel is a Bubble Tea model for a breadcrumb crumbs.
// It manages the crumbs of crumbs and renders them with styling and truncation.
type BreadcrumbsModel struct {
	Styles    Styles
	crumbs    []Crumb
	separator string
	width     int
	row, col  int         // Screen position (set by parent via SetPosition)
	mouse     *mouseState // Shared pointer — survives value copies
}

func (m BreadcrumbsModel) Separator() string {
	return m.separator
}

func (m BreadcrumbsModel) Width() int {
	return m.width
}

// NewBreadcrumbsModel creates a new BreadcrumbsModel with default styles.
func NewBreadcrumbsModel() BreadcrumbsModel {
	return BreadcrumbsModel{
		Styles:    DefaultStyles(),
		separator: " > ",
		mouse:     &mouseState{hoveredIdx: -1},
	}
}

// Push appends a crumb to the crumbs and returns a copy.
func (m BreadcrumbsModel) Push(crumb Crumb) BreadcrumbsModel {
	newCrumbs := make([]Crumb, len(m.crumbs)+1)
	copy(newCrumbs, m.crumbs)
	newCrumbs[len(m.crumbs)] = crumb
	m.crumbs = newCrumbs
	return m
}

// Pop removes the last crumb from the crumbs and returns a copy.
// No-op if the crumbs is empty.
func (m BreadcrumbsModel) Pop() BreadcrumbsModel {
	if len(m.crumbs) == 0 {
		return m
	}
	newCrumbs := make([]Crumb, len(m.crumbs)-1)
	copy(newCrumbs, m.crumbs[:len(m.crumbs)-1])
	m.crumbs = newCrumbs
	return m
}

// SetCrumbs replaces the entire crumbs and returns a copy.
func (m BreadcrumbsModel) SetCrumbs(crumbs []Crumb) BreadcrumbsModel {
	newCrumbs := make([]Crumb, len(crumbs))
	copy(newCrumbs, crumbs)
	m.crumbs = newCrumbs
	return m
}

// SetCrumb updates the crumb at the given index and returns a copy.
// No-op if the index is out of range.
func (m BreadcrumbsModel) SetCrumb(index int, crumb Crumb) BreadcrumbsModel {
	var crumbs []Crumb

	if index < 0 || index >= len(m.crumbs) {
		goto end
	}
	crumbs = make([]Crumb, len(m.crumbs))
	copy(crumbs, m.crumbs)
	crumbs[index] = crumb
	m.crumbs = crumbs
end:
	return m
}

// Crumbs returns a copy of the current crumbs.
func (m BreadcrumbsModel) Crumbs() []Crumb {
	out := make([]Crumb, len(m.crumbs))
	copy(out, m.crumbs)
	return out
}

// Len returns the number of crumbs in the crumbs.
func (m BreadcrumbsModel) Len() int {
	return len(m.crumbs)
}

// SetSize sets the available width for the breadcrumb crumbs.
func (m BreadcrumbsModel) SetSize(width int) BreadcrumbsModel {
	m.width = width
	return m
}

// SetPosition sets the screen position of the breadcrumb crumbs.
// The parent model must call this so that mouse hit testing works correctly.
func (m BreadcrumbsModel) SetPosition(row, col int) BreadcrumbsModel {
	m.row = row
	m.col = col
	return m
}

// Position returns the screen position set by SetPosition.
func (m BreadcrumbsModel) Position() (row, col int) {
	return m.row, m.col
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
//
//goland:noinspection GoAssignmentToReceiver
func (m BreadcrumbsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case PushCrumbMsg:
		m = m.Push(msg.Crumb)
	case PopCrumbMsg:
		m = m.Pop()
	case SetCrumbsMsg:
		m = m.SetCrumbs(msg.Crumbs)
	case SetCrumbMsg:
		m = m.SetCrumb(msg.Index, msg.Crumb)
	}
	return m, nil
}

// View implements tea.Model. Renders the breadcrumb crumbs with mouse support.
func (m BreadcrumbsModel) View() tea.View {
	result := m.render()
	m.mouse.crumbBounds = result.bounds

	v := tea.NewView(result.content)
	v.MouseMode = tea.MouseModeAllMotion
	v.OnMouse = m.HandleMouse
	return v
}

// HitTest returns the crumb index at position (x, y), or -1 if none.
func (m BreadcrumbsModel) HitTest(x, y int) int {
	if y != m.row {
		return -1
	}
	relX := x - m.col
	for i, b := range m.mouse.crumbBounds {
		if relX >= b.startX && relX < b.endX {
			return i
		}
	}
	return -1
}

// HandleMouse processes mouse events from the View's OnMouse callback.
func (m BreadcrumbsModel) HandleMouse(msg tea.MouseMsg) tea.Cmd {
	mouse := msg.Mouse()
	idx := m.HitTest(mouse.X, mouse.Y)

	switch msg.(type) {
	case tea.MouseClickMsg:
		if idx >= 0 && idx < len(m.crumbs) {
			crumb := m.crumbs[idx]
			return func() tea.Msg {
				return CrumbClickedMsg{
					Index:  idx,
					Crumb:  crumb,
					Button: mouse.Button,
				}
			}
		}

	case tea.MouseMotionMsg:
		if idx == m.mouse.hoveredIdx {
			return nil // dedup — same crumb, don't re-emit
		}
		m.mouse.hoveredIdx = idx
		if idx < 0 {
			return func() tea.Msg {
				return CrumbHoverLeaveMsg{}
			}
		}
		if idx < len(m.crumbs) {
			crumb := m.crumbs[idx]
			return func() tea.Msg {
				return CrumbHoverMsg{
					Index: idx,
					Crumb: crumb,
				}
			}
		}
	}

	return nil
}
