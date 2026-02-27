package teatree

// Package-level configuration helpers and presets

type ExpanderControls struct {
	Expand        string
	Collapse      string
	NotApplicable string
}

var PlusExpanderControls = ExpanderControls{
	Expand:        "+",
	Collapse:      "─",
	NotApplicable: "",
}
var NoExpanderControls = ExpanderControls{
	Expand:        "",
	Collapse:      "",
	NotApplicable: "",
}
var TriangleExpanderControls = ExpanderControls{
	Expand:        "▶",
	Collapse:      "▼",
	NotApplicable: "",
}

// BranchStyle defines a set of characters for tree structure renderingq
type BranchStyle struct {
	// Vertical is the character for vertical continuation (e.g., "│ ", "| ")
	Vertical string

	// Horizontal is the character for horizontal continuation (e.g., "│ ", "| ")
	Horizontal string

	// MiddleChild is the character for middle children (e.g., "├─ ", "+- ")
	MiddleChild string

	// LastChild is the character for the last child (e.g., "└─ ", "`- ")
	LastChild string

	// PreExpanderIndent is the space(s) or tab after the prefix
	PreExpanderIndent string

	// PreIconIndent is the space(s) or tab before icon
	PreIconIndent string

	// PreTextIndent is the space(s) or tab before icon
	PreTextIndent string

	// PreSuffixIndent is the space(s) or tab before suffix
	PreSuffixIndent string

	// EmptySpace is the character for empty space (e.g., "  ", "   ")
	EmptySpace string

	ExpanderControls ExpanderControls
}

// Predefined branch styles

// DefaultBranchStyle uses compact spacing with Unicode box-drawing characters
func DefaultBranchStyle(ec ExpanderControls) BranchStyle {
	return BranchStyle{
		PreExpanderIndent: " ", // One space before expander
		PreIconIndent:     " ", // One space before text
		PreSuffixIndent:   " ", // One space before suffix
		PreTextIndent:     " ", // One space before suffix
		ExpanderControls:  ec,
	}
}

// CompactBranchStyle uses compact spacing with Unicode box-drawing characters
func CompactBranchStyle(ec ExpanderControls) (bs BranchStyle) {
	bs = DefaultBranchStyle(ec)
	bs.Vertical = "│"         // Vertical line + space for next level
	bs.Horizontal = "─"       // Horizonal line + space for next level
	bs.MiddleChild = "├─"     // Branch + space before filename
	bs.LastChild = "└─"       // Last branch + space before filename
	bs.EmptySpace = "  "      // Two spaces to match vertical line maxWidth
	bs.PreExpanderIndent = "" // No spaces before expander
	return bs
}

// ASCIIBranchStyle uses ASCII characters for maximum compatibility
func ASCIIBranchStyle(ec ExpanderControls) BranchStyle {
	return BranchStyle{
		Vertical:          "| ",
		MiddleChild:       "+- ",
		LastChild:         "`- ",
		EmptySpace:        "  ",
		PreExpanderIndent: " ", // One space before expander
		PreIconIndent:     " ", // One space before text
		PreSuffixIndent:   " ", // One space before suffix
		ExpanderControls:  ec,
	}
}

// WideBranchStyle uses 4-space indentation (like standard treeview)
func WideBranchStyle(ec ExpanderControls) BranchStyle {
	return BranchStyle{
		Vertical:          "│   ",
		MiddleChild:       "├── ",
		LastChild:         "└── ",
		EmptySpace:        "    ",
		PreExpanderIndent: " ", // One space before expander
		PreIconIndent:     " ", // One space before text
		PreSuffixIndent:   " ", // One space before suffix
		ExpanderControls:  ec,
	}
}

// MinimalBranchStyle uses minimal spacing (1 char per level)
func MinimalBranchStyle(ec ExpanderControls) BranchStyle {
	return BranchStyle{
		Vertical:          "│",
		MiddleChild:       "├",
		LastChild:         "└",
		EmptySpace:        " ",
		PreExpanderIndent: " ", // One space before expander
		PreIconIndent:     " ", // One space before text
		PreSuffixIndent:   " ", // One space before suffix
		ExpanderControls:  ec,
	}
}
