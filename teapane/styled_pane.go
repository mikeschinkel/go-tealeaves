package teapane

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/tealayout"
)

// ContentFunc renders pane content given the content dimensions and focus state.
type ContentFunc func(width, height int, focused bool) string

// StyledPane is a boilerplate-free pane widget that implements all five
// tealayout interfaces: ContentProvider, SetSizer, Styler, SizeHinter, and Focusable.
//
// The caller provides a BorderStyle (for frame appearance) and a ContentFunc
// (for rendering). StyledPane handles size storage, focus state management,
// style rebuilding on focus changes, and SizeHint computation from frame sizes.
type StyledPane struct {
	border      BorderStyle
	content     ContentFunc
	style       lipgloss.Style
	label       string
	flexPercent float64 // percentage of visible flex space (set by app)
	width       int     // content width (frame already subtracted by tealayout)
	height      int     // content height
	focused      bool
	minWidth     int // minimum content width for SizeHint
	sizeHintFunc func(availW, availH int) tealayout.SizeHint
}

// NewStyledPane creates a StyledPane with the given border and content callback.
func NewStyledPane(border BorderStyle, content ContentFunc) *StyledPane {
	p := &StyledPane{
		border:  border,
		content: content,
	}
	p.style = border.Build(false)
	return p
}

// WithMinWidth sets the minimum content width reported in SizeHint.
func (p *StyledPane) WithMinWidth(n int) *StyledPane {
	p.minWidth = n
	return p
}

// WithSizeHintFunc sets a custom SizeHint function. When set, SizeHint()
// delegates to this function instead of computing from minWidth and frame sizes.
// This is useful for panes wrapping content whose desired size is dynamic
// (e.g., a tree whose width depends on its data).
func (p *StyledPane) WithSizeHintFunc(fn func(availW, availH int) tealayout.SizeHint) *StyledPane {
	p.sizeHintFunc = fn
	return p
}

// WithLabel sets a display label for this pane.
func (p *StyledPane) WithLabel(label string) *StyledPane {
	p.label = label
	return p
}

// Label returns the pane's display label.
func (p *StyledPane) Label() string { return p.label }

// FlexPercent returns the pane's percentage of visible flex space.
func (p *StyledPane) FlexPercent() float64 { return p.flexPercent }

// SetFlexPercent sets the pane's percentage of visible flex space.
// Typically called by the app from View() using MultiPaneLayout.VisibleFlexPercents().
func (p *StyledPane) SetFlexPercent(pct float64) { p.flexPercent = pct }

// SetContentFunc replaces the content rendering callback. This is useful when
// the callback needs to capture the StyledPane pointer itself (e.g. to read
// Label or FlexPercent), which isn't possible at NewStyledPane time.
func (p *StyledPane) SetContentFunc(fn ContentFunc) { p.content = fn }

// Width returns the current content width.
func (p *StyledPane) Width() int { return p.width }

// Height returns the current content height.
func (p *StyledPane) Height() int { return p.height }

// Focused returns whether this pane currently has focus.
func (p *StyledPane) Focused() bool { return p.focused }

// SetSize stores the content dimensions. tealayout calls this after
// subtracting the frame (because StyledPane implements Styler).
func (p *StyledPane) SetSize(w, h int) {
	p.width = w
	p.height = h
}

// Style returns the current lipgloss style (used by tealayout for
// frame subtraction).
func (p *StyledPane) Style() lipgloss.Style {
	return p.style
}

// SizeHint reports size preferences. If a custom SizeHintFunc was set via
// WithSizeHintFunc, it is called instead. Otherwise, Min.Width includes
// frame + minWidth and Min.Height includes frame + content line count.
func (p *StyledPane) SizeHint(availW, availH int) tealayout.SizeHint {
	if p.sizeHintFunc != nil {
		return p.sizeHintFunc(availW, availH)
	}
	hint := tealayout.SizeHint{
		Min: tealayout.Size{
			Width: p.border.FrameWidth() + p.minWidth,
		},
	}
	if p.content != nil {
		contentW := p.width
		if contentW <= 0 {
			contentW = availW - p.border.FrameWidth()
		}
		if contentW < 1 {
			contentW = 1
		}
		contentH := p.height
		if contentH <= 0 {
			contentH = availH - p.border.FrameHeight()
		}
		if contentH < 1 {
			contentH = 1
		}
		rendered := p.content(contentW, contentH, p.focused)
		lines := strings.Count(rendered, "\n") + 1
		hint.Min.Height = p.border.FrameHeight() + lines
	}
	return hint
}

// Content renders the pane: content inside the styled frame.
// Implements tealayout.ContentProvider.
func (p *StyledPane) Content() string {
	var rendered string
	if p.content != nil {
		rendered = p.content(p.width, p.height, p.focused)
	}
	totalW := p.width + p.border.FrameWidth()
	totalH := p.height + p.border.FrameHeight()
	return p.style.Width(totalW).Height(totalH).Render(rendered)
}

// Focus sets focused state and rebuilds the style.
func (p *StyledPane) Focus() {
	p.focused = true
	p.style = p.border.Build(true)
}

// Blur clears focused state and rebuilds the style.
func (p *StyledPane) Blur() {
	p.focused = false
	p.style = p.border.Build(false)
}
