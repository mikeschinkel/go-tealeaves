package teadiffview

// SplitPaneRow holds one row of a side-by-side diff display.
// Each row has a left (old) and right (new) pane line.
type SplitPaneRow struct {
	CommitLine PaneLine // Left pane (old/commit version)
	ActualLine PaneLine // Right pane (new/actual version)
	BlockIndex int      // 0 for unchanged, >0 for change blocks (1-based)
	LineOffset int      // 0-based offset within the block (for line-level assignments)
}

// PaneLine is the interface for content displayed in a single pane cell.
type PaneLine interface {
	LineNo() int
	PaneLine()
}

var _ PaneLine = (*TextLine)(nil)

// TextLine holds actual text content with a line number.
type TextLine struct {
	lineNo int
	Text   string
}

func (tl TextLine) LineNo() int { return tl.lineNo }
func (TextLine) PaneLine()     {}

// NewTextLine creates a new TextLine instance.
func NewTextLine(lineNo int, text string) *TextLine {
	return &TextLine{
		lineNo: lineNo,
		Text:   text,
	}
}

var _ PaneLine = (*BlockMarker)(nil)

// BlockMarker indicates the start of a change block on the opposite pane.
// It marks where insertions or deletions begin.
type BlockMarker struct {
	lineNo    int
	LineCount int
}

func (BlockMarker) PaneLine()       {}
func (bm BlockMarker) LineNo() int  { return bm.lineNo }

// NewBlockMarker creates a new BlockMarker instance.
func NewBlockMarker(lineNo int, lineCount int) *BlockMarker {
	return &BlockMarker{
		lineNo:    lineNo,
		LineCount: lineCount,
	}
}

// IsBlockStart returns true (all BlockMarker instances mark block starts).
func (bm *BlockMarker) IsBlockStart() bool {
	return true
}

var _ PaneLine = (*PlaceholderLine)(nil)

// PlaceholderLine is a blank row in one pane that keeps the other pane
// aligned during insertions or deletions.
type PlaceholderLine struct {
	lineNo   int
	HunkLine int
}

func (PlaceholderLine) PaneLine()       {}
func (hr PlaceholderLine) LineNo() int  { return hr.lineNo }

// NewPlaceholderLine creates a new PlaceholderLine instance.
func NewPlaceholderLine(lineNo int, hunkLine int) *PlaceholderLine {
	return &PlaceholderLine{
		lineNo:   lineNo,
		HunkLine: hunkLine,
	}
}

// IsWithinHunk returns true (all PlaceholderLine instances are within hunks).
func (pl *PlaceholderLine) IsWithinHunk() bool {
	return true
}
