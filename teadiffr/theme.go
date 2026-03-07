package teadiffr

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// NewThemedTUIRenderer creates a TUIRenderer with colors derived from the
// given theme's palette. This is an alternative to NewTUIRenderer when using
// the theming system.
func NewThemedTUIRenderer(theme teautils.Theme) *TUIRenderer {
	p := theme.System
	return &TUIRenderer{
		FileHeaderColor:    p.TextMuted,
		BlockHeaderColor:   p.AccentAlt,
		ContextColor:       p.TextDim,
		AddedColor:         p.StatusSuccess,
		DeletedColor:       p.StatusError,
		NewStatusColor:     p.StatusSuccess,
		DeletedStatusColor: p.StatusError,
		NewBgColor:         p.TintPositive,
		DeletedBgColor:     p.TintNegative,
	}
}
