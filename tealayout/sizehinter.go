package tealayout

// SizeHinter is implemented by widgets that can report their size preferences
// given available space. The layout engine calls SizeHint during resolution
// for Fit() constrained children.
type SizeHinter interface {
	SizeHint(availWidth, availHeight int) SizeHint
}
