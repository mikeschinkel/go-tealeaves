package teatree

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Renderer renders a tree to a string
type Renderer[T any] struct {
	tree *Tree[T]
}

// NewRenderer creates a new renderer for the given tree
func NewRenderer[T any](tree *Tree[T]) *Renderer[T] {
	return &Renderer[T]{
		tree: tree,
	}
}

// Render renders the tree to a string
func (r *Renderer[T]) Render() string {
	var lines []string

	for _, root := range r.tree.Nodes() {
		r.renderNode(root, &lines)
	}

	return strings.Join(lines, "\n")
}

// renderNode recursively renders a node and its visible children
func (r *Renderer[T]) renderNode(node *Node[T], lines *[]string) {
	if !node.IsVisible() {
		return
	}

	provider := r.tree.Provider()

	var sb strings.Builder
	branchStyle := provider.BranchStyle()
	expanderControl := provider.ExpanderControl(node)
	icon := provider.Icon(node)
	emptySpace := branchStyle.EmptySpace
	if icon == "" && len(emptySpace) >= 2 {
		emptySpace = emptySpace[:len(emptySpace)-1]
	}
	// Add vertical lines or spaces for each ancestor level
	for _, isLastChild := range node.AncestorIsLastChild() {
		if isLastChild {
			// Ancestor was last child - use empty space
			sb.WriteString(emptySpace)
			s := strings.Repeat(" ", len([]rune(branchStyle.MiddleChild+branchStyle.PreExpanderIndent))-1)
			sb.WriteString(s)
			continue
		}
		// Ancestor has more siblings - use vertical line
		sb.WriteString(branchStyle.Vertical)
		sb.WriteString(emptySpace)
		if icon != "" {
			sb.WriteString(branchStyle.PreIconIndent)
		}
		sb.WriteString(branchStyle.PreExpanderIndent)
	}

	// Add the branch character for this node
	if node.IsLastChild() {
		sb.WriteString(branchStyle.LastChild)
	} else {
		sb.WriteString(branchStyle.MiddleChild)
	}
	hasChildren := node.HasChildren()
	if !hasChildren && len(expanderControl) > 0 {
		sb.WriteString(strings.Repeat(branchStyle.Horizontal, len([]rune(expanderControl+branchStyle.Horizontal))))
	}

	// Build the line components
	//prefix := provider.Prefix(node)
	text := provider.Text(node)
	suffix := provider.Suffix(node)

	if hasChildren {
		indent := provider.BranchStyle().PreExpanderIndent
		if text != "" {
			sb.WriteString(indent)
		}
		sb.WriteString(expanderControl)
	}
	preIconIndent := ""
	if icon != "" {
		preIconIndent = branchStyle.PreIconIndent
	}
	isRoot := node.parent.IsRoot()
	switch {
	case !hasChildren && isRoot:
		fallthrough
	case !isRoot && node.parent.HasGrandChildren() && !node.HasChildren():
		preIconIndent = branchStyle.Horizontal + preIconIndent
	}
	if text != "" && preIconIndent != "" {
		sb.WriteString(preIconIndent)
	}
	sb.WriteString(icon)

	sb.WriteString(branchStyle.PreTextIndent)
	// Apply styling
	style := provider.Style(node, r.tree)
	sb.WriteString(style.Render(text))

	if suffix != "" {
		sb.WriteString(branchStyle.PreSuffixIndent)
		sb.WriteString(suffix)
	}

	*lines = append(*lines, sb.String())

	// Recursively render children if expanded
	if node.IsExpanded() {
		for _, child := range node.Children() {
			r.renderNode(child, lines)
		}
	}
}

// RenderToLines renders the tree to a slice of strings (one per line)
// Useful for viewport integration
func (r *Renderer[T]) RenderToLines() []string {
	var lines []string

	for _, root := range r.tree.Nodes() {
		r.renderNode(root, &lines)
	}

	return lines
}

// GetMaxLineWidth calculates the maximum maxWidth needed to display all visible nodes
// without truncation (ANSI codes stripped for accurate measurement)
func (r *Renderer[T]) GetMaxLineWidth() int {
	lines := r.RenderToLines()
	maxWidth := 0

	for _, line := range lines {
		// Use ansi.StringWidth for accurate display-width measurement
		// (multi-byte tree characters like │ ├ ─ ▶ ▼ are 3 bytes but 1 column)
		width := ansi.StringWidth(line)

		if width > maxWidth {
			maxWidth = width
		}
	}

	return maxWidth
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(s string) string {
	result := strings.Builder{}
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			inEscape = true
			i++ // Skip the '['
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}

	return result.String()
}
