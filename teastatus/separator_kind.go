package teastatus

// SeparatorKind identifies which separator style to use for indicators.
type SeparatorKind int

const (
	// PipeSeparator separates indicators with " | "
	PipeSeparator SeparatorKind = iota
	// SpaceSeparator separates indicators with "  "
	SpaceSeparator
	// BracketSeparator wraps each indicator as "[text]"
	BracketSeparator
)
