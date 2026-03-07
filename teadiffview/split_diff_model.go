package teadiffview

import (
	"fmt"
	"image/color"
	"log/slog"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mikeschinkel/go-diffutils"
)

const scrollStep = 4 // Characters to scroll per horizontal key press

// SelectAction tracks Ctrl+A toggle state.
type SelectAction int

const (
	SelectActionNone SelectAction = iota
	SelectActionBlock
	SelectActionAllBlocks
)

// RowAnnotation marks a row with a visual indicator in the gutter column.
// The model renders the gutter column automatically — consumers just provide data.
type RowAnnotation struct {
	Char  rune        // Gutter character ('A', '●', '!', '→')
	Color color.Color // Color for the character
}

// SplitDiffModelArgs holds configuration for NewSplitDiffModel.
type SplitDiffModelArgs struct {
	Width             int
	Height            int
	Logger            *slog.Logger
	HighlightFunc     func(text, language string) string  // Optional syntax highlighting
	InlineHighlighter diffutils.InlineHighlighter         // Optional char-level highlighting
}

// SplitDiffModel displays two diff panes side-by-side with synchronized scrolling.
// It is a Bubble Tea Model that accepts [diffutils.DiffContent] for display.
type SplitDiffModel struct {
	// Shared state
	xOffset int // Horizontal scroll offset (shared)

	// Cursor and scroll for aligned diff mode
	cursorIndex     int // Current index in rows array
	cursorScreenRow int // Desired screen row for cursor (0 = top, height-1 = bottom)

	// Selection state
	selectionStart   int          // -1 if no selection
	selectionEnd     int          // -1 if no selection (inclusive)
	lastSelectAction SelectAction // Track Ctrl+A toggle state

	// Dimensions
	splitContentWidth int // Width of content for ONE side (excluding line numbers)
	width             int // Total width (both sides + separator + borders)
	height            int

	// Line numbers
	leftLineNumWidth  int // Width of left pane line number column
	rightLineNumWidth int // Width of right pane line number column

	// Gutter state (annotations)
	gutterChars  []rune        // Character per row
	gutterColors []color.Color // Color per row

	// Split pane diff mode
	rows          []SplitPaneRow // Split pane rows
	contentLoaded bool           // True after SetContent called

	// Display
	focused bool

	// Optional highlighting
	highlightFunc     func(text, language string) string
	inlineHighlighter diffutils.InlineHighlighter

	// Logging
	Logger *slog.Logger
}

// NewSplitDiffModel creates a new split diff pane model.
func NewSplitDiffModel(args *SplitDiffModelArgs) SplitDiffModel {
	if args == nil {
		args = &SplitDiffModelArgs{}
	}
	splitContentWidth := args.Width
	if splitContentWidth == 0 {
		splitContentWidth = 40
	}
	height := args.Height
	if height == 0 {
		height = 20
	}
	return SplitDiffModel{
		splitContentWidth: splitContentWidth,
		width:             5 + 2*splitContentWidth,
		height:            height,
		leftLineNumWidth:  4,
		rightLineNumWidth: 4,
		selectionStart:    -1,
		selectionEnd:      -1,
		highlightFunc:     args.HighlightFunc,
		inlineHighlighter: args.InlineHighlighter,
		Logger:            args.Logger,
	}
}

// SetContent loads diff data into the model for display.
func (m SplitDiffModel) SetContent(content *diffutils.DiffContent) SplitDiffModel {
	var maxLeftNum int
	var maxRightNum int

	m.rows = buildSplitPaneDiff(content, m.Logger)
	m.contentLoaded = true

	// Apply syntax highlighting if a highlight function is provided
	if m.highlightFunc != nil && content != nil {
		language := content.Label // Use label as language hint
		for i := range m.rows {
			if textLine, ok := m.rows[i].CommitLine.(*TextLine); ok {
				textLine.Text = m.highlightFunc(textLine.Text, language)
			}
			if textLine, ok := m.rows[i].ActualLine.(*TextLine); ok {
				textLine.Text = m.highlightFunc(textLine.Text, language)
			}
		}
	}

	// Track max line numbers
	for i := range m.rows {
		if m.rows[i].CommitLine != nil && m.rows[i].CommitLine.LineNo() > maxLeftNum {
			maxLeftNum = m.rows[i].CommitLine.LineNo()
		}
		if m.rows[i].ActualLine != nil && m.rows[i].ActualLine.LineNo() > maxRightNum {
			maxRightNum = m.rows[i].ActualLine.LineNo()
		}
	}

	m.leftLineNumWidth = calculateLineNumberWidth(maxLeftNum)
	m.rightLineNumWidth = calculateLineNumberWidth(maxRightNum)

	// Reset state
	m.cursorIndex = 0
	m.cursorScreenRow = 0
	m.xOffset = 0
	m.selectionStart = -1
	m.selectionEnd = -1
	m.lastSelectAction = SelectActionNone

	return m
}

// SetAnnotations applies per-row annotations. The model renders the gutter.
func (m SplitDiffModel) SetAnnotations(annotations map[int]RowAnnotation) SplitDiffModel {
	if len(m.rows) == 0 {
		return m
	}
	chars := make([]rune, len(m.rows))
	colors := make([]color.Color, len(m.rows))
	for i := range chars {
		chars[i] = ' '
	}
	for idx, ann := range annotations {
		if idx >= 0 && idx < len(m.rows) {
			chars[idx] = ann.Char
			colors[idx] = ann.Color
		}
	}
	m.gutterChars = chars
	m.gutterColors = colors
	return m
}

// SetGutter updates gutter indicators directly with parallel arrays.
// chars and colors arrays should have one entry per row in m.rows.
func (m SplitDiffModel) SetGutter(chars []rune, colors []color.Color) SplitDiffModel {
	if len(chars) == len(m.rows) && len(colors) == len(m.rows) {
		m.gutterChars = chars
		m.gutterColors = colors
	}
	return m
}

// calculateLineNumberWidth returns the width needed for line numbers.
func calculateLineNumberWidth(lineCount int) int {
	if lineCount == 0 {
		return 4
	}
	numDigits := len(fmt.Sprintf("%d", lineCount))
	if numDigits == 1 {
		return 4 // "  1 "
	}
	return 1 + numDigits + 1
}

// formatLineNumber formats a line number to a fixed width.
func formatLineNumber(lineNum, width int) string {
	numStr := fmt.Sprintf("%d", lineNum)
	numDigits := len(numStr)
	leadingSpaces := width - numDigits - 1
	if leadingSpaces < 1 {
		leadingSpaces = 1
	}
	return strings.Repeat(" ", leadingSpaces) + numStr + " "
}

// --- Navigation ---

// MoveCursorUp moves cursor up in the aligned lines array.
func (m SplitDiffModel) MoveCursorUp() SplitDiffModel {
	if m.cursorIndex > 0 {
		currentInBlock := m.cursorIndex < len(m.rows) && m.rows[m.cursorIndex].BlockIndex > 0
		m.cursorIndex--
		if !currentInBlock {
			if m.cursorScreenRow > 0 {
				m.cursorScreenRow--
			}
		}
	}
	return m
}

// MoveCursorDown moves cursor down in the aligned lines array.
func (m SplitDiffModel) MoveCursorDown() SplitDiffModel {
	if len(m.rows) == 0 {
		return m
	}
	maxCursorIndex := len(m.rows) - 1 + m.height
	if m.cursorIndex < maxCursorIndex {
		currentInBlock := m.cursorIndex < len(m.rows) && m.rows[m.cursorIndex].BlockIndex > 0
		m.cursorIndex++
		pastEndOfFile := m.cursorIndex >= len(m.rows)
		if !currentInBlock && !pastEndOfFile {
			if m.cursorScreenRow < m.height-1 {
				m.cursorScreenRow++
			}
		}
	}
	return m
}

// PageDown moves cursor down by one page.
func (m SplitDiffModel) PageDown() SplitDiffModel {
	if len(m.rows) == 0 {
		return m
	}
	pageSize := m.height - 1
	if pageSize < 1 {
		pageSize = 1
	}
	targetIndex := m.cursorIndex + pageSize
	maxIndex := len(m.rows) - 1 + m.height
	if targetIndex > maxIndex {
		targetIndex = maxIndex
	}
	m.cursorIndex = targetIndex
	m = m.calculateCursorScreenRow()
	return m
}

// PageUp moves cursor up by one page.
func (m SplitDiffModel) PageUp() SplitDiffModel {
	if len(m.rows) == 0 {
		return m
	}
	pageSize := m.height - 1
	if pageSize < 1 {
		pageSize = 1
	}
	targetIndex := m.cursorIndex - pageSize
	if targetIndex < 0 {
		targetIndex = 0
	}
	m.cursorIndex = targetIndex
	m = m.calculateCursorScreenRow()
	return m
}

// GoToTop jumps to the beginning of the file.
func (m SplitDiffModel) GoToTop() SplitDiffModel {
	m.cursorIndex = 0
	m.cursorScreenRow = 0
	return m
}

// GoToBottom jumps to the end of the file.
func (m SplitDiffModel) GoToBottom() SplitDiffModel {
	if len(m.rows) == 0 {
		return m
	}
	m.cursorIndex = len(m.rows) - 1
	m.cursorScreenRow = m.height - 1
	if m.cursorIndex < m.cursorScreenRow {
		m.cursorScreenRow = m.cursorIndex
	}
	return m
}

// --- Selection ---

// ExtendSelectionUp extends or starts selection upward.
func (m SplitDiffModel) ExtendSelectionUp() SplitDiffModel {
	if m.cursorIndex <= 0 {
		return m
	}
	if m.selectionStart == -1 {
		m.selectionStart = m.cursorIndex
		m.selectionEnd = m.cursorIndex
	}
	m.cursorIndex--
	m.selectionStart = m.cursorIndex
	m.lastSelectAction = SelectActionNone
	return m
}

// ExtendSelectionDown extends or starts selection downward.
func (m SplitDiffModel) ExtendSelectionDown() SplitDiffModel {
	if m.cursorIndex >= len(m.rows)-1 {
		return m
	}
	if m.selectionStart == -1 {
		m.selectionStart = m.cursorIndex
		m.selectionEnd = m.cursorIndex
	}
	m.cursorIndex++
	m.selectionEnd = m.cursorIndex
	m.lastSelectAction = SelectActionNone
	return m
}

// ClearSelection clears any active selection.
func (m SplitDiffModel) ClearSelection() SplitDiffModel {
	m.selectionStart = -1
	m.selectionEnd = -1
	m.lastSelectAction = SelectActionNone
	return m
}

// HasSelection returns true if there is an active selection.
func (m SplitDiffModel) HasSelection() bool {
	return m.selectionStart != -1 && m.selectionEnd != -1
}

// SelectCurrentBlock selects all rows in the current change block.
func (m SplitDiffModel) SelectCurrentBlock() SplitDiffModel {
	start, end := m.getBlockRange(m.cursorIndex)
	if start == -1 || end == -1 {
		return m
	}
	if m.cursorIndex < len(m.rows) && m.rows[m.cursorIndex].BlockIndex == 0 {
		m.selectionStart = m.cursorIndex
		m.selectionEnd = m.cursorIndex
		return m
	}
	m.selectionStart = start
	m.selectionEnd = end
	return m
}

// SelectAllBlocks selects all change blocks in the file.
func (m SplitDiffModel) SelectAllBlocks() SplitDiffModel {
	if len(m.rows) == 0 {
		return m
	}
	firstBlock := -1
	lastBlock := -1
	for i, row := range m.rows {
		if row.BlockIndex > 0 {
			if firstBlock == -1 {
				firstBlock = i
			}
			lastBlock = i
		}
	}
	if firstBlock == -1 {
		return m
	}
	m.selectionStart = firstBlock
	m.selectionEnd = lastBlock
	m.cursorIndex = firstBlock
	m.cursorScreenRow = 0
	return m
}

// ToggleBlockSelection cycles through: current block -> all blocks -> no selection.
func (m SplitDiffModel) ToggleBlockSelection() SplitDiffModel {
	switch m.lastSelectAction {
	case SelectActionNone:
		m = m.SelectCurrentBlock()
		m.lastSelectAction = SelectActionBlock
	case SelectActionBlock:
		m = m.SelectAllBlocks()
		m.lastSelectAction = SelectActionAllBlocks
	case SelectActionAllBlocks:
		m = m.ClearSelection()
		m.lastSelectAction = SelectActionNone
	}
	return m
}

// IsRowSelected returns true if the given row index is within the selection range.
func (m SplitDiffModel) IsRowSelected(rowIdx int) bool {
	if m.selectionStart == -1 || m.selectionEnd == -1 {
		return false
	}
	return rowIdx >= m.selectionStart && rowIdx <= m.selectionEnd
}

// --- Horizontal scrolling ---

// ScrollLeft scrolls horizontally to the left.
func (m SplitDiffModel) ScrollLeft() SplitDiffModel {
	m.xOffset -= scrollStep
	if m.xOffset < 0 {
		m.xOffset = 0
	}
	return m
}

// ScrollRight scrolls horizontally to the right.
func (m SplitDiffModel) ScrollRight() SplitDiffModel {
	m.xOffset += scrollStep
	return m
}

// ScrollToColumn scrolls to a specific column.
func (m SplitDiffModel) ScrollToColumn(col int) SplitDiffModel {
	m.xOffset = col
	if m.xOffset < 0 {
		m.xOffset = 0
	}
	return m
}

// ScrollToEnd scrolls to show the end of the longest line.
func (m SplitDiffModel) ScrollToEnd() SplitDiffModel {
	var maxWidth int
	for _, row := range m.rows {
		if textLine, ok := row.CommitLine.(*TextLine); ok {
			w := ansi.StringWidth(textLine.Text)
			if w > maxWidth {
				maxWidth = w
			}
		}
		if textLine, ok := row.ActualLine.(*TextLine); ok {
			w := ansi.StringWidth(textLine.Text)
			if w > maxWidth {
				maxWidth = w
			}
		}
	}
	if maxWidth > m.splitContentWidth {
		m.xOffset = maxWidth - m.splitContentWidth
	}
	return m
}

// --- Size and focus ---

// SetSize updates dimensions.
func (m SplitDiffModel) SetSize(splitContentWidth, height int) SplitDiffModel {
	m.splitContentWidth = splitContentWidth
	m.width = 5 + 2*splitContentWidth
	m.height = height
	return m
}

// Focus sets the model to focused state.
func (m SplitDiffModel) Focus() SplitDiffModel {
	m.focused = true
	return m
}

// Blur sets the model to unfocused state.
func (m SplitDiffModel) Blur() SplitDiffModel {
	m.focused = false
	return m
}

// --- Accessors ---

// LeftLineNumWidth returns the width of the left pane line number column.
func (m SplitDiffModel) LeftLineNumWidth() int {
	return m.leftLineNumWidth
}

// RightLineNumWidth returns the width of the right pane line number column.
func (m SplitDiffModel) RightLineNumWidth() int {
	return m.rightLineNumWidth
}

// GetSelectedLines returns the selected line range (start, end inclusive).
// Returns (-1, -1) if no selection is active.
func (m SplitDiffModel) GetSelectedLines() (start, end int) {
	return m.selectionStart, m.selectionEnd
}

// LineCount returns the number of rows in the diff.
func (m SplitDiffModel) LineCount() int {
	return len(m.rows)
}

// RowCount returns the number of rows in the diff.
func (m SplitDiffModel) RowCount() int {
	return len(m.rows)
}

// Rows returns the split pane rows.
func (m SplitDiffModel) Rows() []SplitPaneRow {
	return m.rows
}

// CursorIndex returns the current cursor row index.
func (m SplitDiffModel) CursorIndex() int {
	return m.cursorIndex
}

// GetBlockRange returns the start and end row indices of the change block
// containing the cursor. Returns (-1, -1) if cursor is out of bounds.
func (m SplitDiffModel) GetBlockRange() (start, end int) {
	return m.getBlockRange(m.cursorIndex)
}

// GetBlockIndexAtCursor returns the BlockIndex (1-based) of the row at cursor position.
// Returns 0 if cursor is on an unchanged line or out of bounds.
func (m SplitDiffModel) GetBlockIndexAtCursor() int {
	if m.cursorIndex < 0 || m.cursorIndex >= len(m.rows) {
		return 0
	}
	return m.rows[m.cursorIndex].BlockIndex
}

// CenterOnFirstChangeBlock positions the view to show the first change block optimally.
func (m SplitDiffModel) CenterOnFirstChangeBlock() SplitDiffModel {
	if len(m.rows) == 0 {
		return m
	}

	firstChangeIdx := -1
	for i, row := range m.rows {
		if row.BlockIndex > 0 {
			firstChangeIdx = i
			break
		}
	}
	if firstChangeIdx == -1 {
		return m
	}

	firstBlockID := m.rows[firstChangeIdx].BlockIndex
	lastChangeIdx := firstChangeIdx
	for i := firstChangeIdx + 1; i < len(m.rows); i++ {
		if m.rows[i].BlockIndex == firstBlockID {
			lastChangeIdx = i
		} else {
			break
		}
	}

	blockSize := lastChangeIdx - firstChangeIdx + 1

	if firstChangeIdx <= 1 {
		m.cursorIndex = firstChangeIdx
		m.cursorScreenRow = firstChangeIdx
		return m
	}

	if blockSize <= m.height {
		contextAbove := (m.height - blockSize) / 2
		if contextAbove > firstChangeIdx {
			contextAbove = firstChangeIdx
		}
		m.cursorIndex = firstChangeIdx
		m.cursorScreenRow = contextAbove
		return m
	}

	const contextLines = 2
	m.cursorIndex = firstChangeIdx
	if firstChangeIdx >= contextLines {
		m.cursorScreenRow = contextLines
	} else {
		m.cursorScreenRow = firstChangeIdx
	}
	return m
}

// --- tea.Model interface ---

// Init initializes the model.
func (m SplitDiffModel) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m SplitDiffModel) Update(msg tea.Msg) (SplitDiffModel, tea.Cmd) {
	return m, nil
}

// View renders the split diff pane content (without border — parent applies border).
func (m SplitDiffModel) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}
	if len(m.rows) == 0 {
		if m.contentLoaded {
			return tea.NewView("(empty file)")
		}
		return tea.NewView("No diff to display")
	}
	return tea.NewView(m.viewSplitPane())
}

// --- Internal rendering ---

func (m SplitDiffModel) calculateCursorScreenRow() SplitDiffModel {
	if m.cursorIndex >= len(m.rows) {
		if m.cursorScreenRow > m.height-1 {
			m.cursorScreenRow = m.height - 1
		}
		return m
	}

	leftVisible := m.countVisibleLinesUpTo(m.cursorIndex, true)
	rightVisible := m.countVisibleLinesUpTo(m.cursorIndex, false)

	maxScreenRow := leftVisible
	if rightVisible < maxScreenRow {
		maxScreenRow = rightVisible
	}
	if m.cursorScreenRow > maxScreenRow {
		m.cursorScreenRow = maxScreenRow
	}
	if m.cursorScreenRow > m.height-1 {
		m.cursorScreenRow = m.height - 1
	}
	if m.cursorScreenRow < 0 {
		m.cursorScreenRow = 0
	}
	return m
}

func (m SplitDiffModel) countVisibleLinesUpTo(targetIdx int, isLeftPane bool) int {
	count := 0
	for i := 0; i < targetIdx && i < len(m.rows); i++ {
		row := m.rows[i]
		if isLeftPane {
			if _, isPlaceholder := row.CommitLine.(*PlaceholderLine); !isPlaceholder {
				count++
			}
		} else {
			if _, isPlaceholder := row.ActualLine.(*PlaceholderLine); !isPlaceholder {
				count++
			}
		}
	}
	return count
}

func (m SplitDiffModel) getBlockRange(cursorRow int) (start, end int) {
	if cursorRow < 0 || cursorRow >= len(m.rows) {
		return -1, -1
	}
	blockID := m.rows[cursorRow].BlockIndex
	if blockID == 0 {
		return cursorRow, cursorRow
	}
	start = cursorRow
	for start > 0 && m.rows[start-1].BlockIndex == blockID {
		start--
	}
	end = cursorRow
	for end < len(m.rows)-1 && m.rows[end+1].BlockIndex == blockID {
		end++
	}
	return start, end
}

func (m SplitDiffModel) calculatePaneStartIdx(cursorIdx, targetScreenRow int, isLeftPane bool) int {
	effectiveCursorIdx := cursorIdx
	if effectiveCursorIdx < len(m.rows) {
		for effectiveCursorIdx > 0 {
			row := m.rows[effectiveCursorIdx]
			var isPlaceholder bool
			if isLeftPane {
				_, isPlaceholder = row.CommitLine.(*PlaceholderLine)
			} else {
				_, isPlaceholder = row.ActualLine.(*PlaceholderLine)
			}
			if !isPlaceholder {
				break
			}
			effectiveCursorIdx--
		}
	}

	if targetScreenRow == 0 {
		if effectiveCursorIdx >= len(m.rows) {
			return len(m.rows) - 1
		}
		return effectiveCursorIdx
	}

	linesToCount := targetScreenRow
	if cursorIdx >= len(m.rows) {
		rowsPastEnd := cursorIdx - (len(m.rows) - 1)
		linesToCount = targetScreenRow - rowsPastEnd
		if linesToCount < 0 {
			linesToCount = 0
		}
	}

	countFromIdx := effectiveCursorIdx
	if countFromIdx >= len(m.rows) {
		countFromIdx = len(m.rows) - 1
	}

	if linesToCount == 0 {
		return countFromIdx
	}

	visibleCount := 0
	idx := countFromIdx - 1

	for idx >= 0 && visibleCount < linesToCount {
		row := m.rows[idx]
		if isLeftPane {
			if _, isPlaceholder := row.CommitLine.(*PlaceholderLine); !isPlaceholder {
				visibleCount++
				if visibleCount == linesToCount {
					return idx
				}
			}
		} else {
			if _, isPlaceholder := row.ActualLine.(*PlaceholderLine); !isPlaceholder {
				visibleCount++
				if visibleCount == linesToCount {
					return idx
				}
			}
		}
		idx--
	}

	return 0
}

func (m SplitDiffModel) viewSplitPane() string {
	if len(m.rows) == 0 {
		if m.contentLoaded {
			return "(empty file)"
		}
		return "No diff to display"
	}

	leftIdx := m.calculatePaneStartIdx(m.cursorIndex, m.cursorScreenRow, true)
	rightIdx := m.calculatePaneStartIdx(m.cursorIndex, m.cursorScreenRow, false)

	var leftSide []string
	var rightSide []string

	for len(leftSide) < m.height && leftIdx < len(m.rows) {
		row := m.rows[leftIdx]
		rowIdx := leftIdx
		isCursor := m.focused && rowIdx == m.cursorIndex
		isSelected := m.focused && m.IsRowSelected(rowIdx) && !isCursor
		isBlockRow := row.BlockIndex > 0

		if _, isPlaceholder := row.CommitLine.(*PlaceholderLine); !isPlaceholder {
			shouldHighlight := isCursor
			if !shouldHighlight && m.focused && isBlockRow {
				if _, isBlockMarker := row.CommitLine.(*BlockMarker); isBlockMarker {
					if m.cursorIndex >= 0 && m.cursorIndex < len(m.rows) {
						cursorRow := m.rows[m.cursorIndex]
						shouldHighlight = (cursorRow.BlockIndex == row.BlockIndex)
					}
				}
			}
			leftSide = append(leftSide, m.renderPaneLine(row.CommitLine, rowIdx, shouldHighlight, isSelected, isBlockRow, true))
		}
		leftIdx++
	}

	for len(rightSide) < m.height && rightIdx < len(m.rows) {
		row := m.rows[rightIdx]
		rowIdx := rightIdx
		isCursor := m.focused && rowIdx == m.cursorIndex
		isSelected := m.focused && m.IsRowSelected(rowIdx) && !isCursor
		isBlockRow := row.BlockIndex > 0

		if _, isPlaceholder := row.ActualLine.(*PlaceholderLine); !isPlaceholder {
			shouldHighlight := isCursor
			if !shouldHighlight && m.focused && isBlockRow {
				if _, isBlockMarker := row.ActualLine.(*BlockMarker); isBlockMarker {
					if m.cursorIndex >= 0 && m.cursorIndex < len(m.rows) {
						cursorRow := m.rows[m.cursorIndex]
						shouldHighlight = (cursorRow.BlockIndex == row.BlockIndex)
					}
				}
			}
			rightSide = append(rightSide, m.renderPaneLine(row.ActualLine, rowIdx, shouldHighlight, isSelected, isBlockRow, false))
		}
		rightIdx++
	}

	// Pad with blank lines
	cursorPastEnd := m.cursorIndex >= len(m.rows)

	for len(leftSide) < m.height {
		screenRow := len(leftSide)
		shouldHighlight := cursorPastEnd && m.focused && screenRow == m.cursorScreenRow
		totalWidth := m.leftLineNumWidth + m.splitContentWidth

		if shouldHighlight {
			blank := strings.Repeat(" ", totalWidth)
			leftSide = append(leftSide, lipgloss.NewStyle().Reverse(true).Render(blank))
		} else {
			leftSide = append(leftSide, strings.Repeat(" ", totalWidth))
		}
	}

	for len(rightSide) < m.height {
		screenRow := len(rightSide)
		shouldHighlight := cursorPastEnd && m.focused && screenRow == m.cursorScreenRow
		totalWidth := m.rightLineNumWidth + m.splitContentWidth

		if shouldHighlight {
			blank := strings.Repeat(" ", totalWidth)
			rightSide = append(rightSide, lipgloss.NewStyle().Reverse(true).Render(blank))
		} else {
			rightSide = append(rightSide, strings.Repeat(" ", totalWidth))
		}
	}

	var result strings.Builder
	for i := 0; i < len(leftSide); i++ {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(leftSide[i])
		result.WriteString("│")
		result.WriteString(rightSide[i])
	}

	return result.String()
}

func (m SplitDiffModel) renderPaneLine(pl PaneLine, rowIdx int, isCursor bool, isSelected bool, isBlockRow bool, isLeft bool) string {
	const gutterWidth = 2
	var lineNumStr string
	var content string
	var lineNumWidth int
	var gutterStr string
	var gutterChar rune
	var gutterColor color.Color
	var totalWidth int
	var isContentSide bool

	if isLeft {
		lineNumWidth = m.leftLineNumWidth
	} else {
		lineNumWidth = m.rightLineNumWidth
	}

	totalWidth = gutterWidth + lineNumWidth + (m.splitContentWidth - gutterWidth)

	_, isPlaceholder := pl.(*PlaceholderLine)
	_, isBlockMarker := pl.(*BlockMarker)
	isContentSide = !isPlaceholder && !isBlockMarker

	gutterChar = ' '
	if isContentSide && rowIdx >= 0 && rowIdx < len(m.gutterChars) {
		gutterChar = m.gutterChars[rowIdx]
	}
	if rowIdx >= 0 && rowIdx < len(m.gutterColors) {
		gutterColor = m.gutterColors[rowIdx]
	}

	if isPlaceholder {
		blankSpace := strings.Repeat(" ", totalWidth)
		if isCursor {
			return lipgloss.NewStyle().Reverse(true).Render(blankSpace)
		}
		if isSelected {
			return lipgloss.NewStyle().Background(lipgloss.Color(SelectionBgColor)).Render(blankSpace)
		}
		if isBlockRow {
			return lipgloss.NewStyle().Background(lipgloss.Color(ChangeBlockBgColor)).Render(blankSpace)
		}
		return blankSpace
	}

	if isBlockMarker {
		blankSpace := strings.Repeat(" ", totalWidth)
		if isCursor {
			return lipgloss.NewStyle().Reverse(true).Render(blankSpace)
		}
		if isSelected {
			return lipgloss.NewStyle().Background(lipgloss.Color(SelectionBgColor)).Render(blankSpace)
		}
		return lipgloss.NewStyle().Background(lipgloss.Color(ChangeBlockBgColor)).Render(blankSpace)
	}

	textLine := pl.(*TextLine)
	if textLine.LineNo() > 0 {
		lineNumStr = formatLineNumber(textLine.LineNo(), lineNumWidth)
		content = textLine.Text
	} else {
		lineNumStr = strings.Repeat(" ", lineNumWidth)
		content = ""
	}

	content = strings.ReplaceAll(content, "\t", "  ")

	if isCursor {
		gutterStr = string(gutterChar) + " "
	} else if gutterChar != ' ' {
		gutterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(gutterColor)
		gutterStr = gutterStyle.Render(string(gutterChar)) + " "
	} else {
		gutterStr = "  "
	}

	if m.xOffset > 0 && content != "" {
		content = ansi.TruncateLeft(content, m.xOffset, "")
	}

	contentWidth := m.splitContentWidth - gutterWidth
	if contentWidth < 0 {
		contentWidth = 0
	}
	content = ansi.Truncate(content, contentWidth, "")

	visualWidth := ansi.StringWidth(content)
	if visualWidth < contentWidth {
		content = content + strings.Repeat(" ", contentWidth-visualWidth)
	}

	if isCursor {
		content = ansi.Strip(content)
		fullLine := gutterStr + lineNumStr + content
		return lipgloss.NewStyle().Reverse(true).Render(fullLine)
	}

	if isSelected {
		content = strings.ReplaceAll(content, ANSIReset, ANSIReset+SelectionBgANSI)
		styledContent := SelectionBgANSI + lineNumStr + content + ANSIReset
		return gutterStr + styledContent
	}

	if isBlockRow {
		content = strings.ReplaceAll(content, ANSIReset, ANSIReset+ChangeBlockBgANSI)
		styledContent := ChangeBlockBgANSI + lineNumStr + content + ANSIReset
		return gutterStr + styledContent
	}

	return gutterStr + lineNumStr + content
}
