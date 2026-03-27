package teatree

import (
	"image/color"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// TreeModel is the BubbleTea model for the tree.
//
// Deprecated name: This type was previously named Model. The deprecated
// [NewModel] constructor is provided for backward compatibility.
type TreeModel[T any] struct {
	Keys     TreeKeyMap // Keyboard bindings
	tree     *Tree[T]
	renderer *Renderer[T]
	viewport viewport.Model
	width    int
	height   int
	ready    bool
	theme    *teautils.Theme

	// Optional self-managed frame (border + padding).
	// When active, SetSize(w, h) = total dimensions INCLUDING frame chrome.
	// The tree internally computes content area, truncates lines, and renders
	// with the frame — ensuring no mismatch between truncation and border width.
	frameStyle  lipgloss.Style
	hasFrame    bool
	header      string         // optional header text rendered inside frame, above tree content
	headerStyle lipgloss.Style // style for header text
}

// NewTreeModel creates a new BubbleTea model for the tree
func NewTreeModel[T any](tree *Tree[T], height int) TreeModel[T] {
	renderer := NewRenderer(tree)
	width := renderer.GetMaxLineWidth()
	return TreeModel[T]{
		Keys:     DefaultTreeKeyMap(),
		tree:     tree,
		renderer: renderer,
		viewport: viewport.New(viewport.WithWidth(width), viewport.WithHeight(height)),
		height:   height,
		ready:    true,
	}
}

// Init implements tea.Model
func (m TreeModel[T]) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m TreeModel[T]) Update(msg tea.Msg) (TreeModel[T], tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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
		m.viewport.SetWidth(m.contentWidth())
		m.viewport.SetHeight(m.treeHeight())
		m = m.updateViewportContent()
		return m, nil
	}

	// Delegate to viewport for scrolling
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

// View implements tea.Model (FOLLOWS ClearPath)
func (m TreeModel[T]) View() (view tea.View) {
	var lines []string
	var start, end int
	var visibleLines []string
	var frameWidth int
	var maxVisibleWidth int
	var w int
	var hw int
	var rw int
	var cw int
	var content string
	var h string

	if !m.ready {
		view = tea.NewView("Initializing...")
		goto end
	}

	// Render content without horizontal padding (viewport pads to maxWidth)
	// We want tree to be only as wide as actual content
	lines = m.renderer.RenderToLines()

	// Apply vertical scrolling from viewport (YOffset)
	start = m.viewport.YOffset()
	end = start + m.treeHeight()

	if end < 0 {
		view = tea.NewView("")
		goto end
	}
	if start >= len(lines) {
		view = tea.NewView("")
		goto end
	}
	if end > len(lines) {
		end = len(lines)
	}

	visibleLines = lines[start:end]

	// Compute frame width from VIEWPORT-VISIBLE lines, not all expanded lines.
	// GetMaxLineWidth() measures ALL expanded nodes, but off-screen nodes with
	// longer names inflate the width, creating a visible gap between tree content
	// and the border. Using only viewport-visible lines keeps the border tight.
	frameWidth = m.width
	if m.hasFrame {
		for _, line := range visibleLines {
			w = ansi.StringWidth(line)
			if w > maxVisibleWidth {
				maxVisibleWidth = w
			}
		}
		// Also consider header width so border is never narrower than header
		if m.header != "" {
			hw = ansi.StringWidth(m.headerStyle.Render(m.header))
			if hw > maxVisibleWidth {
				maxVisibleWidth = hw
			}
		}
		rw = maxVisibleWidth + m.frameStyle.GetHorizontalBorderSize() + m.frameStyle.GetHorizontalPadding()
		if rw < frameWidth || frameWidth == 0 {
			frameWidth = rw
		}
	}

	// Inner content width for truncation (derived from frameWidth, not cached m.width)
	if m.hasFrame {
		cw = frameWidth - m.frameStyle.GetHorizontalBorderSize() - m.frameStyle.GetHorizontalPadding()
	}
	if !m.hasFrame {
		cw = frameWidth
	}

	// Truncate lines to inner content width (prevents wrapping inside frame)
	if cw > 0 {
		for i, line := range visibleLines {
			if ansi.StringWidth(line) > cw {
				visibleLines[i] = ansi.Truncate(line, cw, "…")
			}
		}
	}

	content = joinLines(visibleLines)

	// Prepend header if set, truncated to content width
	if m.header != "" {
		h = m.headerStyle.Render(m.header)
		if cw > 0 && ansi.StringWidth(h) > cw {
			h = ansi.Truncate(h, cw, "…")
		}
		content = h + "\n" + content
	}

	// Wrap in frame if active
	if m.hasFrame {
		content = m.frameStyle.
			Width(frameWidth).
			Height(m.height).
			Render(content)
	}

	view = tea.NewView(content)

end:
	return view
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
func (m TreeModel[T]) updateViewportContent() TreeModel[T] {
	m.viewport.SetContent(m.renderer.Render())
	return m
}

// ensureFocusedVisible scrolls the viewport to ensure the focused node is visible
// TODO Should this be pushed down to teatree.TreeModel instead of being here?
func (m TreeModel[T]) ensureFocusedVisible() TreeModel[T] {
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
	if focusedIndex < m.viewport.YOffset() {
		m.viewport.SetYOffset(focusedIndex)
	}

	// If focused line is below viewport, scroll down
	if focusedIndex >= m.viewport.YOffset()+m.viewport.Height() {
		m.viewport.SetYOffset(focusedIndex - m.viewport.Height() + 1)
	}
	return m
}

// Tree returns the underlying tree
func (m TreeModel[T]) Tree() *Tree[T] {
	return m.tree
}

// SetSize updates the model dimensions.
// When a frame is active, width and height are total rendered dimensions
// including frame chrome. The viewport is sized to the inner content area.
func (m TreeModel[T]) SetSize(width, height int) TreeModel[T] {
	m.width = width
	m.height = height
	m.viewport.SetWidth(m.contentWidth())
	m.viewport.SetHeight(m.treeHeight())
	m = m.updateViewportContent()
	return m
}

// MaxLineWidth returns the maximum line maxWidth needed to display all content
func (m TreeModel[T]) MaxLineWidth() int {
	return m.renderer.GetMaxLineWidth()
}

// FocusedNode returns the currently focused node
func (m TreeModel[T]) FocusedNode() (node *Node[T]) {
	return m.tree.FocusedNode()
}

// SetFocusedNode sets the focused node by ID
func (m TreeModel[T]) SetFocusedNode(nodeID string) TreeModel[T] {
	m.tree.SetFocusedNode(nodeID)
	return m.ensureFocusedVisible()
}

// WithFrame sets a lipgloss style for a self-managed border+padding frame.
// When active, SetSize(w, h) = total rendered dimensions including frame chrome.
// The tree truncates content to the inner width and renders the frame itself.
func (m TreeModel[T]) WithFrame(style lipgloss.Style) TreeModel[T] {
	m.frameStyle = style
	m.hasFrame = true
	return m
}

// WithHeader sets an optional header line rendered inside the frame, above tree content.
func (m TreeModel[T]) WithHeader(text string, style lipgloss.Style) TreeModel[T] {
	m.header = text
	m.headerStyle = style
	return m
}

// SetBorderColor updates the frame border foreground color at runtime (e.g. focus changes).
func (m TreeModel[T]) SetBorderColor(c color.Color) TreeModel[T] {
	m.frameStyle = m.frameStyle.BorderForeground(c)
	return m
}

// contentWidth returns the width available for tree content inside the frame.
func (m TreeModel[T]) contentWidth() int {
	if !m.hasFrame {
		return m.width
	}
	return m.width - m.frameStyle.GetHorizontalBorderSize() - m.frameStyle.GetHorizontalPadding()
}

// contentHeight returns the height available for content inside the frame (includes header).
func (m TreeModel[T]) contentHeight() int {
	if !m.hasFrame {
		return m.height
	}
	return m.height - m.frameStyle.GetVerticalBorderSize() - m.frameStyle.GetVerticalPadding()
}

// treeHeight returns the height available for tree lines (content height minus header).
func (m TreeModel[T]) treeHeight() int {
	h := m.contentHeight()
	if m.header != "" {
		h-- // header takes 1 line
	}
	return h
}

// RequiredWidth returns the total width needed: max content line width + frame chrome.
// Consumers use this for layout calculations.
func (m TreeModel[T]) RequiredWidth() int {
	w := m.renderer.GetMaxLineWidth()
	if m.hasFrame {
		w += m.frameStyle.GetHorizontalBorderSize() + m.frameStyle.GetHorizontalPadding()
	}
	return w
}
