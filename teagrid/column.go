package teagrid

import "charm.land/lipgloss/v2"

const (
	defaultPaddingLeft  = 1
	defaultPaddingRight = 0
)

// Column defines a column in the grid.
type Column struct {
	title        string
	key          string
	width        int
	flexFactor   int
	paddingLeft  int
	paddingRight int
	alignment    lipgloss.Position
	filterable   bool
	style        lipgloss.Style
	fmtString    string
}

// NewColumn creates a new fixed-width column.
// Default alignment is Left, paddingLeft is 1, paddingRight is 0.
func NewColumn(key, title string, width int) Column {
	return Column{
		key:          key,
		title:        title,
		width:        width,
		paddingLeft:  defaultPaddingLeft,
		paddingRight: defaultPaddingRight,
		alignment:    lipgloss.Left,
	}
}

// NewFlexColumn creates a flexible-width column that fills available space.
// Multiple flex columns share space proportional to their flex factors.
func NewFlexColumn(key, title string, flexFactor int) Column {
	return Column{
		key:          key,
		title:        title,
		flexFactor:   max(flexFactor, 1),
		paddingLeft:  defaultPaddingLeft,
		paddingRight: defaultPaddingRight,
		alignment:    lipgloss.Left,
	}
}

// WithStyle applies a style to the column.
func (c Column) WithStyle(style lipgloss.Style) Column {
	c.style = style
	return c
}

// WithFiltered sets whether the column participates in filtering.
func (c Column) WithFiltered(filterable bool) Column {
	c.filterable = filterable
	return c
}

// WithFormatString sets the format string used by fmt.Sprintf to display data.
// Unlike v0.1.0, this applies to both data cells and header cells.
func (c Column) WithFormatString(fmtString string) Column {
	c.fmtString = fmtString
	return c
}

// WithPadding sets left and right cell padding (in characters).
func (c Column) WithPadding(left, right int) Column {
	c.paddingLeft = left
	c.paddingRight = right
	return c
}

// WithPaddingLeft sets the left cell padding.
func (c Column) WithPaddingLeft(padding int) Column {
	c.paddingLeft = padding
	return c
}

// WithPaddingRight sets the right cell padding.
func (c Column) WithPaddingRight(padding int) Column {
	c.paddingRight = padding
	return c
}

// WithAlignment sets the text alignment within the column.
func (c Column) WithAlignment(alignment lipgloss.Position) Column {
	c.alignment = alignment
	return c
}

// RenderWidth returns the total rendered width of the column including padding.
func (c Column) RenderWidth() int {
	return c.paddingLeft + c.contentWidth() + c.paddingRight
}

// contentWidth returns the width available for content (either explicit or
// computed from flex).
func (c Column) contentWidth() int {
	return c.width
}

// Title returns the column title.
func (c Column) Title() string { return c.title }

// Key returns the column key used to match row data.
func (c Column) Key() string { return c.key }

// Width returns the content width of the column (excluding padding).
func (c Column) Width() int { return c.width }

// FlexFactor returns the flex factor, or 0 for fixed columns.
func (c Column) FlexFactor() int { return c.flexFactor }

// IsFlex returns whether this is a flex-width column.
func (c Column) IsFlex() bool { return c.flexFactor != 0 }

// Filterable returns whether the column participates in filtering.
func (c Column) Filterable() bool { return c.filterable }

// Style returns the column style.
func (c Column) Style() lipgloss.Style { return c.style }

// FmtString returns the format string.
func (c Column) FmtString() string { return c.fmtString }

// PaddingLeft returns the left padding.
func (c Column) PaddingLeft() int { return c.paddingLeft }

// PaddingRight returns the right padding.
func (c Column) PaddingRight() int { return c.paddingRight }

// Alignment returns the text alignment.
func (c Column) Alignment() lipgloss.Position { return c.alignment }
