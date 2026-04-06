package tealayout

import "charm.land/lipgloss/v2"

// SetSizer is implemented by widgets that accept size assignments from the
// layout engine. The layout engine calls SetSize with content dimensions
// (total assigned minus border/padding if a style is attached).
type SetSizer interface {
	SetSize(width, height int)
}

// ContentProvider is implemented by widgets that can render themselves to a
// string for use in layout panes. This is distinct from tea.Model.View()
// which returns tea.View — a struct with metadata beyond the content string.
type ContentProvider interface {
	Content() string
}

// Styler is implemented by widgets that expose a lipgloss Style. The layout
// engine uses this to compute content dimensions from total assigned space.
type Styler interface {
	Style() lipgloss.Style
}

// Focusable is implemented by widgets that can receive and lose focus.
type Focusable interface {
	Focus()
	Blur()
}
