package teagrid

import "charm.land/lipgloss/v2"

// BorderRegion controls the appearance of one border region.
type BorderRegion struct {
	Visible bool
	Style   lipgloss.Style
}

// BorderConfig defines the border appearance for each region of the grid.
type BorderConfig struct {
	// Outer controls the outer frame of the grid.
	Outer BorderRegion

	// Header controls the separator between header and data rows.
	Header BorderRegion

	// Inner controls column dividers between cells.
	Inner BorderRegion

	// Footer controls the separator between data rows and footer.
	Footer BorderRegion

	// Chars defines the box-drawing characters used.
	Chars BorderChars
}

// BorderChars defines the characters used for box drawing.
type BorderChars struct {
	Horizontal string
	Vertical   string

	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string

	LeftJunction   string
	RightJunction  string
	TopJunction    string
	BottomJunction string

	InnerJunction string
	InnerDivider  string
}

// Pre-defined character sets.
var (
	charsDefault = BorderChars{
		Horizontal:     "━",
		Vertical:       "┃",
		TopLeft:        "┏",
		TopRight:       "┓",
		BottomLeft:     "┗",
		BottomRight:    "┛",
		LeftJunction:   "┣",
		RightJunction:  "┫",
		TopJunction:    "┳",
		BottomJunction: "┻",
		InnerJunction:  "╋",
		InnerDivider:   "┃",
	}

	charsRounded = BorderChars{
		Horizontal:     "─",
		Vertical:       "│",
		TopLeft:        "╭",
		TopRight:       "╮",
		BottomLeft:     "╰",
		BottomRight:    "╯",
		LeftJunction:   "├",
		RightJunction:  "┤",
		TopJunction:    "┬",
		BottomJunction: "┴",
		InnerJunction:  "┼",
		InnerDivider:   "│",
	}

	charsMinimal = BorderChars{
		Horizontal:   "─",
		InnerDivider: " ",
	}
)

// BorderDefault returns a heavy-weight border configuration.
func BorderDefault() BorderConfig {
	return BorderConfig{
		Outer:  BorderRegion{Visible: true},
		Header: BorderRegion{Visible: true},
		Inner:  BorderRegion{Visible: true},
		Footer: BorderRegion{Visible: true},
		Chars:  charsDefault,
	}
}

// BorderRounded returns a thin, rounded border configuration.
func BorderRounded() BorderConfig {
	return BorderConfig{
		Outer:  BorderRegion{Visible: true},
		Header: BorderRegion{Visible: true},
		Inner:  BorderRegion{Visible: true},
		Footer: BorderRegion{Visible: true},
		Chars:  charsRounded,
	}
}

// Borderless returns a configuration with no borders.
func Borderless() BorderConfig {
	return BorderConfig{
		Outer:  BorderRegion{Visible: false},
		Header: BorderRegion{Visible: false},
		Inner:  BorderRegion{Visible: false},
		Footer: BorderRegion{Visible: false},
	}
}

// BorderMinimal returns a border with only a horizontal header separator.
func BorderMinimal() BorderConfig {
	return BorderConfig{
		Outer:  BorderRegion{Visible: false},
		Header: BorderRegion{Visible: true},
		Inner:  BorderRegion{Visible: false},
		Footer: BorderRegion{Visible: false},
		Chars:  charsMinimal,
	}
}

// WithOuter returns a copy with the outer border region set.
func (bc BorderConfig) WithOuter(region BorderRegion) BorderConfig {
	bc.Outer = region
	return bc
}

// WithHeader returns a copy with the header border region set.
func (bc BorderConfig) WithHeader(region BorderRegion) BorderConfig {
	bc.Header = region
	return bc
}

// WithInner returns a copy with the inner border region set.
func (bc BorderConfig) WithInner(region BorderRegion) BorderConfig {
	bc.Inner = region
	return bc
}

// WithFooter returns a copy with the footer border region set.
func (bc BorderConfig) WithFooter(region BorderRegion) BorderConfig {
	bc.Footer = region
	return bc
}

// WithChars returns a copy with the specified border characters.
func (bc BorderConfig) WithChars(chars BorderChars) BorderConfig {
	bc.Chars = chars
	return bc
}

// WithStyle returns a copy of the region with the given style.
func (br BorderRegion) WithStyle(style lipgloss.Style) BorderRegion {
	br.Style = style
	return br
}

// WithVisible returns a copy of the region with the given visibility.
func (br BorderRegion) WithVisible(visible bool) BorderRegion {
	br.Visible = visible
	return br
}

// HasOuterBorder returns whether the outer border is visible.
func (bc BorderConfig) HasOuterBorder() bool {
	return bc.Outer.Visible
}

// HasHeaderSeparator returns whether the header separator is visible.
func (bc BorderConfig) HasHeaderSeparator() bool {
	return bc.Header.Visible
}

// HasInnerDividers returns whether column dividers are visible.
func (bc BorderConfig) HasInnerDividers() bool {
	return bc.Inner.Visible
}

// HasFooterSeparator returns whether the footer separator is visible.
func (bc BorderConfig) HasFooterSeparator() bool {
	return bc.Footer.Visible
}

// OuterWidth returns the horizontal width consumed by the outer border (0 or 2).
func (bc BorderConfig) OuterWidth() int {
	if bc.Outer.Visible {
		return 2
	}
	return 0
}

// InnerDividerWidth returns the width of a single inner column divider (0 or 1).
func (bc BorderConfig) InnerDividerWidth() int {
	if bc.Inner.Visible {
		return 1
	}
	return 0
}
