package tealayout

import (
	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teacrumbs"
)

// StackView is the interface for views managed by StackLayoutModel.
// Each view must be a tea.Model and additionally support lifecycle hooks,
// breadcrumb generation, and dimension propagation.
type StackView interface {
	tea.Model
	OnEnter() tea.Cmd
	OnExit() tea.Cmd
	Breadcrumb() teacrumbs.Crumb
	SetSize(width, height int)
}

// PushViewMsg requests pushing a new view onto the stack.
type PushViewMsg struct {
	View     StackView
	CacheKey string
}

// PopViewMsg requests popping the current view.
type PopViewMsg struct{}

// StackLayoutModel is a full tea.Model that manages a stack of views with
// breadcrumbs, push/pop lifecycle, and view caching. It uses pointer-semantics
// fields (slices, map) so mutations survive Bubble Tea's value-copy Update cycle.
type StackLayoutModel struct {
	stack            []StackView
	cacheKeys        []string             // parallel to stack
	cache            map[string]StackView // viewer cache
	crumbs           teacrumbs.BreadcrumbsModel
	width            int
	height           int
	breadcrumbHeight int // height reserved for breadcrumbs (default 1)
}

// NewStackLayoutModel creates a StackLayoutModel with an initial view.
func NewStackLayoutModel(initial StackView, styles teacrumbs.Styles) StackLayoutModel {
	crumbs := teacrumbs.NewBreadcrumbsModel().WithStyles(styles)
	crumbs = crumbs.Push(initial.Breadcrumb())
	return StackLayoutModel{
		stack:            []StackView{initial},
		cacheKeys:        []string{""},
		cache:            make(map[string]StackView),
		crumbs:           crumbs,
		breadcrumbHeight: 1,
	}
}

// Push pushes a view onto the stack and calls lifecycle hooks.
// Returns a tea.Cmd that combines OnExit of the current view and OnEnter of
// the new view. If cacheKey is non-empty, the view is stored in the cache.
func (m *StackLayoutModel) Push(view StackView, cacheKey string) tea.Cmd {
	var cmds []tea.Cmd

	// OnExit for current
	if cur := m.current(); cur != nil {
		if cmd := cur.OnExit(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Push the new view
	m.stack = append(m.stack, view)
	m.cacheKeys = append(m.cacheKeys, cacheKey)

	// Cache if keyed
	if cacheKey != "" {
		m.cache[cacheKey] = view
	}

	// Set size and breadcrumb
	view.SetSize(m.width, m.viewHeight())
	m.crumbs = m.crumbs.Push(view.Breadcrumb())

	// OnEnter for new view
	if cmd := view.OnEnter(); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// Pop removes the current view and returns to the previous one.
// Returns a combined tea.Cmd for lifecycle hooks.
func (m *StackLayoutModel) Pop() (tea.Cmd, error) {
	if len(m.stack) == 0 {
		return nil, NewErr(ErrStackEmpty)
	}
	if len(m.stack) <= 1 {
		return nil, NewErr(ErrStackUnderflow, "depth", len(m.stack))
	}
	return m.popOne(), nil
}

// PopTo pops views until the stack has the given depth.
// depth must be >= 1.
func (m *StackLayoutModel) PopTo(depth int) (tea.Cmd, error) {
	if depth < 1 {
		return nil, NewErr(ErrStackUnderflow, "target_depth", depth)
	}
	if depth > len(m.stack) {
		return nil, NewErr(ErrStackUnderflow, "target_depth", depth, "current_depth", len(m.stack))
	}

	var cmds []tea.Cmd
	for len(m.stack) > depth {
		if cmd := m.popOne(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...), nil
}

// popOne pops a single view off the stack and returns lifecycle cmds.
func (m *StackLayoutModel) popOne() tea.Cmd {
	var cmds []tea.Cmd

	// OnExit for current
	top := m.stack[len(m.stack)-1]
	if cmd := top.OnExit(); cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Pop
	m.stack = m.stack[:len(m.stack)-1]
	m.cacheKeys = m.cacheKeys[:len(m.cacheKeys)-1]
	m.crumbs = m.crumbs.Pop()

	// OnEnter for new current
	if cur := m.current(); cur != nil {
		cur.SetSize(m.width, m.viewHeight())
		if cmd := cur.OnEnter(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// Current returns the view at the top of the stack.
func (m StackLayoutModel) Current() StackView {
	return m.current()
}

func (m StackLayoutModel) current() StackView {
	if len(m.stack) == 0 {
		return nil
	}
	return m.stack[len(m.stack)-1]
}

// Depth returns the number of views on the stack.
func (m StackLayoutModel) Depth() int {
	return len(m.stack)
}

// CanPop returns true if the stack has more than one view.
func (m StackLayoutModel) CanPop() bool {
	return len(m.stack) > 1
}

// GetCached returns a cached view by key.
func (m StackLayoutModel) GetCached(key string) (StackView, bool) {
	v, ok := m.cache[key]
	return v, ok
}

// DeleteCached removes a view from the cache.
func (m *StackLayoutModel) DeleteCached(key string) {
	delete(m.cache, key)
}

// UpdateCurrent replaces the current top-of-stack view.
// This is used after tea.Model.Update returns a new model value.
func (m *StackLayoutModel) UpdateCurrent(updated StackView) {
	if len(m.stack) > 0 {
		m.stack[len(m.stack)-1] = updated

		// Update cache entry if keyed
		key := m.cacheKeys[len(m.cacheKeys)-1]
		if key != "" {
			m.cache[key] = updated
		}
	}
}

// SetSize updates the StackLayoutModel's dimensions. Call this when the host
// application handles WindowSizeMsg directly instead of delegating to
// StackLayoutModel.Update(). This ensures popOne() and Push() propagate
// correct dimensions to views via their SetSize methods.
func (m *StackLayoutModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// BreadcrumbsView renders just the breadcrumb trail for manual composition.
// Use this when the host application needs to compose breadcrumbs, view content,
// and other chrome (status bar, overlays) separately rather than using View().
func (m StackLayoutModel) BreadcrumbsView(width int) tea.View {
	return m.crumbs.SetSize(width).View()
}

// viewHeight returns the height available for the current view.
func (m StackLayoutModel) viewHeight() int {
	h := m.height - m.breadcrumbHeight
	if h < 0 {
		h = 0
	}
	return h
}

// --- tea.Model implementation ---

// Init implements tea.Model.
func (m StackLayoutModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	if cur := m.current(); cur != nil {
		if cmd := cur.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
		if cmd := cur.OnEnter(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// Update implements tea.Model.
//
//goland:noinspection GoAssignmentToReceiver
func (m StackLayoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.crumbs = m.crumbs.SetSize(m.width)
		if cur := m.current(); cur != nil {
			cur.SetSize(m.width, m.viewHeight())
		}
		return m, nil

	case PushViewMsg:
		cmd := m.Push(msg.View, msg.CacheKey)
		return m, cmd

	case PopViewMsg:
		cmd, _ := m.Pop()
		return m, cmd
	}

	// Delegate to current view
	cur := m.current()
	if cur == nil {
		return m, nil
	}
	updated, cmd := cur.Update(msg)
	if sv, ok := updated.(StackView); ok {
		m.UpdateCurrent(sv)
	}
	return m, cmd
}

// View implements tea.Model.
func (m StackLayoutModel) View() tea.View {
	crumbView := m.crumbs.View()

	cur := m.current()
	if cur == nil {
		return crumbView
	}

	curView := cur.View()

	// Compose: breadcrumbs on top, current view below.
	// Merge Content strings with a newline; preserve the current view's
	// mouse handler, cursor, and other tea.View fields.
	content := crumbView.Content + "\n" + curView.Content
	curView.SetContent(content)

	// Preserve the crumbs' mouse handler by wrapping it with the current
	// view's handler. The crumbs handler takes priority for its row.
	if crumbView.OnMouse != nil {
		origMouse := curView.OnMouse
		crumbMouse := crumbView.OnMouse
		curView.OnMouse = func(msg tea.MouseMsg) tea.Cmd {
			// Let crumbs handle it first
			if cmd := crumbMouse(msg); cmd != nil {
				return cmd
			}
			if origMouse != nil {
				return origMouse(msg)
			}
			return nil
		}
	}

	return curView
}
