package teapane

import lipgloss "charm.land/lipgloss/v2"

// PlainPane is a simple borderless pane implementing ContentProvider and SetSizer.
// Use for headers, footers, status bars, and other non-interactive elements.
type PlainPane struct {
	content  ContentFunc
	style    lipgloss.Style
	hasStyle bool
	width    int
	height   int
}

// NewPlainPane creates a PlainPane with the given content callback.
// The focused parameter in ContentFunc is always false for PlainPane.
func NewPlainPane(content ContentFunc) *PlainPane {
	return &PlainPane{
		content: content,
	}
}

// WithStyle sets a lipgloss style for rendering (e.g. background color).
func (p *PlainPane) WithStyle(s lipgloss.Style) *PlainPane {
	p.style = s
	p.hasStyle = true
	return p
}

// Width returns the current width.
func (p *PlainPane) Width() int { return p.width }

// Height returns the current height.
func (p *PlainPane) Height() int { return p.height }

// SetSize stores the dimensions.
func (p *PlainPane) SetSize(w, _ int) {
	p.width = w
}

// Content renders the pane content with the optional style applied.
// Implements tealayout.ContentProvider.
func (p *PlainPane) Content() string {
	var rendered string
	if p.content != nil {
		rendered = p.content(p.width, p.height, false)
	}
	if !p.hasStyle {
		return rendered
	}
	return p.style.Width(p.width).Render(rendered)
}
