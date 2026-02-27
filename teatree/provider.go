package teatree

import (
	"github.com/charmbracelet/lipgloss"
)

// NodeProvider defines the interface for customizing node rendering
type NodeProvider[T any] interface {
	// Icon returns the leading glyph (folder expand/collapse indicator, file marker, etc.)
	Icon(node *Node[T]) string

	// Text returns the formatted display name for the node
	Text(node *Node[T]) string

	// Suffix returns the text to display after Text
	Suffix(node *Node[T]) string

	// Style returns the lipgloss style for the node based on focus state
	//Style(node *Node[T], isFocused bool) lipgloss.Style
	Style(node *Node[T], tree *Tree[T]) lipgloss.Style

	ExpanderControl(node *Node[T]) string

	BranchStyle() BranchStyle
}

// CompactNodeProvider is a default provider with configurable branch style
type CompactNodeProvider[T any] struct {
	// branchStyle defines the branch characters and spacing
	branchStyle BranchStyle
}

func (p *CompactNodeProvider[T]) BranchStyle() BranchStyle {
	return p.branchStyle
}

// NewCompactNodeProvider creates a new compact provider using CompactStyle defaults
func NewCompactNodeProvider[T any](ec ExpanderControls) *CompactNodeProvider[T] {
	return &CompactNodeProvider[T]{
		branchStyle: CompactBranchStyle(ec),
	}
}

func (p *CompactNodeProvider[T]) ExpanderControl(node *Node[T]) string {
	ecs := p.branchStyle.ExpanderControls
	switch {
	case len(node.children) == 0:
		return ecs.NotApplicable
	case node.expanded:
		return ecs.Collapse
	default:
		//fallthrough
	}
	return ecs.Expand
}

// Icon returns a space character (compact providers typically don't add extra icons)
// Override this method to add custom icons (e.g., folder/file indicators)
func (p *CompactNodeProvider[T]) Icon(node *Node[T]) string {
	//if node.HasChildren() {
	//	return "📁"
	//}
	//return "📄"
	return ""
}

func (p *CompactNodeProvider[T]) Suffix(node *Node[T]) string {
	return ""
}

// Text returns the node's name by default
// Override this method to customize the display text
func (p *CompactNodeProvider[T]) Text(node *Node[T]) string {
	return node.Name()
}

// Style returns a default style (foreground color based on focus)
// Override this method to customize colors and styling
func (p *CompactNodeProvider[T]) Style(node *Node[T], tree *Tree[T]) lipgloss.Style {
	if tree.IsFocusedNode(node) {
		return lipgloss.NewStyle().Reverse(true)
	}
	return lipgloss.NewStyle()
}

// SimpleNodeProvider is a minimal provider with no tree structure characters
// Useful for flat lists or when you want complete custom control
type SimpleNodeProvider[T any] struct{}

// Icon returns empty string
func (p *SimpleNodeProvider[T]) Icon(node *Node[T]) string {
	return ""
}

// Format returns the node's name
func (p *SimpleNodeProvider[T]) Format(node *Node[T]) string {
	return node.Name()
}

// Style returns a default style
func (p *SimpleNodeProvider[T]) Style(node *Node[T], isFocused bool) lipgloss.Style {
	if isFocused {
		return lipgloss.NewStyle().Bold(true)
	}
	return lipgloss.NewStyle()
}
