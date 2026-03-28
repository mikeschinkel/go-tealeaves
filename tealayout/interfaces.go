package tealayout

import "charm.land/lipgloss/v2"

// SetSizer is implemented by widgets that accept size assignments from the
// layout engine. The layout engine calls SetSize with content dimensions
// (total assigned minus border/padding if a style is attached).
type SetSizer interface {
	SetSize(width, height int)
}

// Viewer is implemented by widgets that can render themselves to a string.
type Viewer interface {
	View() string
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
