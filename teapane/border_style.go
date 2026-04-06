package teapane

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// BorderStyle configures the visual frame around a pane: border type,
// colors (normal and focused), text foreground, and padding.
type BorderStyle struct {
	Border       lipgloss.Border
	Color        color.Color // border color (normal)
	FocusedColor color.Color // border color when focused (nil = same as Color)
	Foreground   color.Color // text foreground (nil = terminal default)
	PaddingH     int         // horizontal padding (default 1)
	PaddingV     int         // vertical padding (default 0)
}

// Build returns a lipgloss.Style configured from the BorderStyle.
// When focused is true and FocusedColor is set, it uses FocusedColor
// for the border; otherwise it uses Color.
func (bs BorderStyle) Build(focused bool) lipgloss.Style {
	borderColor := bs.Color
	if focused && bs.FocusedColor != nil {
		borderColor = bs.FocusedColor
	}

	s := lipgloss.NewStyle().
		Border(bs.Border).
		Padding(bs.PaddingV, bs.PaddingH)

	if borderColor != nil {
		s = s.BorderForeground(borderColor)
	}
	if bs.Foreground != nil {
		s = s.Foreground(bs.Foreground)
	}
	return s
}

// FrameWidth returns the total horizontal space consumed by border + padding.
func (bs BorderStyle) FrameWidth() int {
	s := lipgloss.NewStyle().
		Border(bs.Border).
		Padding(bs.PaddingV, bs.PaddingH)
	return s.GetHorizontalBorderSize() + s.GetHorizontalPadding()
}

// FrameHeight returns the total vertical space consumed by border + padding.
func (bs BorderStyle) FrameHeight() int {
	s := lipgloss.NewStyle().
		Border(bs.Border).
		Padding(bs.PaddingV, bs.PaddingH)
	return s.GetVerticalBorderSize() + s.GetVerticalPadding()
}
