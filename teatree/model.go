package teatree

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the BubbleTea model for the tree
type Model[T any] struct {
	Keys     TreeKeyMap // Keyboard bindings
	tree     *Tree[T]
	renderer *Renderer[T]
	viewport viewport.Model
	width    int
	height   int
	ready    bool
}

// NewModel creates a new BubbleTea model for the tree
func NewModel[T any](tree *Tree[T], height int) Model[T] {
	renderer := NewRenderer(tree)
	width := renderer.GetMaxLineWidth()
	return Model[T]{
		Keys:     DefaultTreeKeyMap(),
		tree:     tree,
		renderer: renderer,
		viewport: viewport.New(width, height),
		height:   height,
		ready:    true,
	}
}

// Init implements tea.Model
func (m Model[T]) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Up):
			if m.tree.MoveUp() {
				m = m.ensureFocusedVisible()
			}
			return m, nil

		case key.Matches(msg, m.Keys.Down):
			if m.tree.MoveDown() {
				m = m.ensureFocusedVisible()
			}
			return m, nil

		case key.Matches(msg, m.Keys.ExpandOrEnter):
			if m.tree.ExpandFocused() {
				// Expanded - update viewport content
				return m.updateViewportContent(), nil
			}
			focused := m.tree.FocusedNode()
			if focused.HasChildren() && focused.IsExpanded() {
				// Already expanded - move to first child
				if m.tree.MoveDown() {
					m = m.ensureFocusedVisible()
				}
			}
			return m, nil

		case key.Matches(msg, m.Keys.CollapseOrUp):
			focused := m.tree.FocusedNode()
			if focused != nil && focused.HasChildren() && focused.IsExpanded() {
				// Collapse if expanded
				m.tree.CollapseFocused()
				m = m.updateViewportContent()
			} else if focused != nil && focused.Parent() != nil {
				// Move to parent if collapsed or no children
				m.tree.SetFocusedNode(focused.Parent().ID())
				m = m.ensureFocusedVisible()
			}
			return m, nil

		case key.Matches(msg, m.Keys.Toggle):
			// Toggle expansion
			if m.tree.ToggleFocused() {
				m = m.updateViewportContent()
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m = m.updateViewportContent()
		return m, nil
	}

	// Delegate to viewport for scrolling
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

// View implements tea.Model
func (m Model[T]) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Render content without horizontal padding (viewport pads to maxWidth)
	// We want tree to be only as wide as actual content
	lines := m.renderer.RenderToLines()

	// Apply vertical scrolling from viewport (YOffset)
	start := m.viewport.YOffset
	end := start + m.viewport.Height

	if end < 0 {
		return ""
	}
	if start >= len(lines) {
		return ""
	}
	if end > len(lines) {
		end = len(lines)
	}

	visibleLines := lines[start:end]

	// TODO: Add horizontal scrolling support for deep paths if needed
	// (viewport.XOffset is unexported, so we'd need to track it ourselves)

	return joinLines(visibleLines)
}

// joinLines joins lines with newlines, handling empty slices
func joinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, line := range lines {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(line)
	}
	return sb.String()
}

// updateViewportContent updates the viewport with the current tree rendering
func (m Model[T]) updateViewportContent() Model[T] {
	m.viewport.SetContent(m.renderer.Render())
	return m
}

// ensureFocusedVisible scrolls the viewport to ensure the focused node is visible
// TODO Should this be pushed down to teatree.Model instead of being here?
func (m Model[T]) ensureFocusedVisible() Model[T] {
	m = m.updateViewportContent()

	// Find the line index of the focused node
	focused := m.tree.FocusedNode()
	if focused == nil {
		return m
	}

	visibleNodes := m.tree.VisibleNodes()
	focusedIndex := -1
	for i, node := range visibleNodes {
		if node == focused {
			focusedIndex = i
			break
		}
	}

	if focusedIndex < 0 {
		return m
	}

	// Scroll viewport to show focused line
	// If focused line is above viewport, scroll up
	if focusedIndex < m.viewport.YOffset {
		m.viewport.YOffset = focusedIndex
	}

	// If focused line is below viewport, scroll down
	if focusedIndex >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.YOffset = focusedIndex - m.viewport.Height + 1
	}
	return m
}

// Tree returns the underlying tree
func (m Model[T]) Tree() *Tree[T] {
	return m.tree
}

// SetSize updates the model dimensions
func (m Model[T]) SetSize(width, height int) Model[T] {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height
	m = m.updateViewportContent()
	return m
}

// MaxLineWidth returns the maximum line maxWidth needed to display all content
func (m Model[T]) MaxLineWidth() int {
	return m.renderer.GetMaxLineWidth()
}

// FocusedNode returns the currently focused node
func (m Model[T]) FocusedNode() (node *Node[T]) {
	return m.tree.FocusedNode()
}

// SetFocusedNode sets the focused node by ID
func (m Model[T]) SetFocusedNode(nodeID string) Model[T] {
	m.tree.SetFocusedNode(nodeID)
	return m.ensureFocusedVisible()
}
