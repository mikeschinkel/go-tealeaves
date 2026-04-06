package teapane

// ScrollContentFunc renders scrollable content given dimensions and scroll offset.
type ScrollContentFunc func(width, height, offset int) string

// ScrollPane extends StyledPane with scroll state management.
// It tracks the scroll offset and total line count, clamping the offset
// to valid bounds automatically.
type ScrollPane struct {
	*StyledPane
	scrollContent ScrollContentFunc
	offset        int
	totalLines    int
}

// NewScrollPane creates a ScrollPane with the given border and scroll content callback.
func NewScrollPane(border BorderStyle, content ScrollContentFunc) *ScrollPane {
	sp := &ScrollPane{
		scrollContent: content,
	}
	// Create the inner StyledPane with a content func that delegates to scrollContent.
	sp.StyledPane = NewStyledPane(border, func(w, h int, focused bool) string {
		return sp.scrollContent(w, h, sp.offset)
	})
	return sp
}

// ScrollOffset returns the current scroll offset.
func (sp *ScrollPane) ScrollOffset() int { return sp.offset }

// SetTotalLines sets the total number of content lines for clamping.
func (sp *ScrollPane) SetTotalLines(n int) {
	sp.totalLines = n
	sp.clampOffset()
}

// ScrollUp moves the viewport up by one line.
func (sp *ScrollPane) ScrollUp() {
	if sp.offset > 0 {
		sp.offset--
	}
}

// ScrollDown moves the viewport down by one line.
func (sp *ScrollPane) ScrollDown() {
	sp.offset++
	sp.clampOffset()
}

// clampOffset ensures offset stays within valid bounds.
func (sp *ScrollPane) clampOffset() {
	maxOffset := sp.totalLines - sp.height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if sp.offset > maxOffset {
		sp.offset = maxOffset
	}
	if sp.offset < 0 {
		sp.offset = 0
	}
}
