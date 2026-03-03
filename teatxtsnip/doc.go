// Package teatxtsnip wraps Bubble Tea's textarea with text selection and
// clipboard support (Ctrl+C/X/V, Shift+Arrow).
//
// It can operate in multi-line or single-line mode via [TextSnipModelArgs].
//
// Usage:
//
//	editor := teatxtsnip.NewTextSnipModel(nil) // multi-line
//	single := teatxtsnip.NewTextSnipModel(&teatxtsnip.TextSnipModelArgs{
//	    SingleLine: true,
//	})
//
// # Stability
//
// This package is provisional as of v0.3.0. The public API may change in
// minor releases until promoted to stable.
package teatxtsnip
