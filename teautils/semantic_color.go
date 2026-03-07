package teautils

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// SemanticColor wraps a color.Color with pre-built lipgloss.Style values for
// foreground, background, and border use. This avoids allocating a new Style on
// every render frame — callers use .Foreground() instead of
// lipgloss.NewStyle().Foreground(c).
//
// SemanticColor implements color.Color so it can be passed directly to lipgloss
// methods that accept color.Color.
type SemanticColor struct {
	raw color.Color
	fg  lipgloss.Style // pre-built: NewStyle().Foreground(raw)
	bg  lipgloss.Style // pre-built: NewStyle().Background(raw)
	bdr lipgloss.Style // pre-built: NewStyle().BorderForeground(raw)
}

// NewSemanticColor creates a SemanticColor with pre-built cached styles.
// Passing nil is safe — lipgloss treats nil color as "no color".
func NewSemanticColor(c color.Color) SemanticColor {
	base := lipgloss.NewStyle()
	return SemanticColor{
		raw: c,
		fg:  base.Foreground(c),
		bg:  base.Background(c),
		bdr: base.BorderForeground(c),
	}
}

// RGBA implements color.Color. Returns (0,0,0,0) if the underlying color is nil.
func (s SemanticColor) RGBA() (r, g, b, a uint32) {
	if s.raw == nil {
		return 0, 0, 0, 0
	}
	return s.raw.RGBA()
}

// Color returns the underlying color.Color (may be nil).
func (s SemanticColor) Color() color.Color {
	return s.raw
}

// IsZero returns true if the underlying color is nil.
func (s SemanticColor) IsZero() bool {
	return s.raw == nil
}

// Foreground returns a pre-built style with this color as foreground.
func (s SemanticColor) Foreground() lipgloss.Style {
	return s.fg
}

// Background returns a pre-built style with this color as background.
func (s SemanticColor) Background() lipgloss.Style {
	return s.bg
}

// BorderForeground returns a pre-built style with this color as border foreground.
func (s SemanticColor) BorderForeground() lipgloss.Style {
	return s.bdr
}

// Render is shorthand for s.Foreground().Render(text).
func (s SemanticColor) Render(text string) string {
	return s.fg.Render(text)
}
